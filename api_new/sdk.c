#include "sdk.h"

const int cacheSize = 4194304; //4mb

/*
** The structure of a plugin handle looks like this:
**
** (struct PluginSharedMemory)
**     toolId (int)
**     toolConfig (void *)
**     engine (struct EngineInterface*)
**     outputAnchors (struct OutputAnchor*)
**         name (void *)
**         isOpen (uint32_t)
**         firstChild (struct OutputConn*)
**             isOpen (uint32_t)
**             ii (struct IncomingInterface*)
**             nextConnection (struct OutputConn*)
**         nextAnchor (struct OutputAnchor*)
**         recordCache (char)
**     totalInputAnchors (uint32_t)
**     closedInputAnchors (uint32_t)
**     inputAnchors (struct InputAnchor*)
**         name (void *)
**         type (void *)
**         isOpen (uint32_t)
**         totalConnections (uint32_t)
**         closedConnections (uint32_t)
**         firstChild (struct InputConnection*)
**             isOpen (uint32_t)
**             metadata (void *)
**             percent (double)
**             nextConnection (struct InputConnection*)
**             plugin (struct PluginSharedMemory*)
**             recordCache (char[])
**         nextAnchor (struct InputAnchor*)
*/

struct InputConnection {
    uint32_t                   isOpen;
    void*                      metadata;
    double                     percent;
    struct InputConnection*    nextConnection;
    struct PluginSharedMemory* plugin;
    char                       recordCache[4194304];
};

struct InputAnchor {
    void*                   name;
    void*                   type;
    uint32_t                isOpen;
    uint32_t                totalConnections;
    uint32_t                closedConnections;
    struct InputConnection* firstChild;
    struct InputAnchor*     nextAnchor;
};

struct OutputConn {
    uint32_t                            isOpen;
    struct IncomingConnectionInterface* ii;
    struct OutputConn*                  nextConnection;
};

struct OutputAnchor {
    void*                name;
    void*                metadata;
    uint32_t             isOpen;
    struct OutputConn*   firstChild;
    struct OutputAnchor* nextAnchor;
    char                 recordCache[4194304];
};

struct PluginSharedMemory {
    uint32_t                toolId;
    void*                   toolConfig;
    struct EngineInterface* engine;
    struct OutputAnchor*    outputAnchors;
    uint32_t                totalInputs;
    uint32_t                closedInputs;
    struct InputAnchor*     inputAnchors;
};

long configurePlugin(int nToolID, void * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface) {
    struct PluginSharedMemory* plugin = malloc(sizeof(struct PluginSharedMemory));
    plugin->toolId = nToolID;
    plugin->toolConfig = pXmlProperties;
    plugin->engine = pEngineInterface;
    plugin->outputAnchors = NULL;
    plugin->totalInputs = 0;
    plugin->closedInputs = 0;
    plugin->inputAnchors = NULL;

    r_pluginInterface->handle = plugin;
    r_pluginInterface->pPI_Close = &PI_Close;
    r_pluginInterface->pPI_PushAllRecords = &PI_PushAllRecords;
    r_pluginInterface->pPI_AddIncomingConnection = &PI_AddIncomingConnection;
    r_pluginInterface->pPI_AddOutgoingConnection = &PI_AddOutgoingConnection;

    return 1;
}

void PI_Close(void * handle, bool bHasErrors) {
    // do nothing
}

long PI_PushAllRecords(void * handle, __int64 nRecordLimit){
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    //goOnComplete(plugin->toolId, nRecordLimit);
}

struct InputAnchor* getOrCreateAnchor(struct PluginSharedMemory* plugin, const wchar_t* name) {
    return NULL;
}

long PI_AddIncomingConnection(void * handle, void * pIncomingConnectionType, void * pIncomingConnectionName, struct IncomingConnectionInterface *r_IncConnInt) {
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    struct InputAnchor *anchor = getOrCreateAnchor(plugin, (const wchar_t*)pIncomingConnectionName);
    struct InputConnection *connection = malloc(sizeof(struct InputConnection));
    connection->isOpen = 1;
    connection->metadata = NULL;
    connection->percent = 0;
    connection->nextConnection = NULL;
    connection->plugin = plugin;

    anchor->totalConnections++;

    /*
    r_IncConnInt->handle = input;
    r_IncConnInt->pII_Init = &II_Init;
    r_IncConnInt->pII_PushRecord = &II_PushRecord;
    r_IncConnInt->pII_UpdateProgress = &II_UpdateProgress;
    r_IncConnInt->pII_Close = &II_Close;
    r_IncConnInt->pII_Free = &II_Free;
    */

    return 1;
}

struct OutputAnchor* getOutputAnchorByName(struct OutputAnchor* anchor, void* name) {
    if (NULL == anchor) {
        return NULL;
    }
    if (wcscmp((const wchar_t*)name, (const wchar_t*)anchor->name) == 0) {
        return anchor;
    }
    return getOutputAnchorByName(anchor->nextAnchor, name);
}

void appendOutgoingConnection(struct OutputAnchor* anchor, struct IncomingConnectionInterface* ii) {
    struct OutputConn* conn = malloc(sizeof(struct OutputConn));
    conn->isOpen = 1;
    conn->ii = ii;
    conn->nextConnection = NULL;

    if (NULL == anchor->firstChild) {
        anchor->firstChild = conn;
        return;
    }

    struct OutputConn *childConn = anchor->firstChild;
    while (childConn->nextConnection != NULL) {
        childConn = childConn->nextConnection;
    }
    childConn->nextConnection = conn;
    if (anchor->isOpen == 1) {
        long result = ii->pII_Init(ii->handle, anchor->metadata);
        if (result == 0) {
            conn->isOpen = 0;
        }
    }
}

struct OutputAnchor* appendAnchor(struct PluginSharedMemory* plugin, void * name) {
    struct OutputAnchor* anchor = malloc(sizeof(struct OutputAnchor));
    anchor->name = name;
    anchor->metadata = NULL;
    anchor->isOpen = 1;
    anchor->firstChild = NULL;
    anchor->nextAnchor = NULL;

    if (NULL == plugin->outputAnchors) {
        plugin->outputAnchors = anchor;
        return anchor;
    }

    struct OutputAnchor* child = plugin->outputAnchors;
    while (NULL != child) {
        child = child->nextAnchor;
    }
    child->nextAnchor = anchor;
    return anchor;
}

long PI_AddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt) {
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    struct OutputAnchor* anchor = getOutputAnchorByName(plugin->outputAnchors, pOutgoingConnectionName);
    if (NULL == anchor) {
        anchor = appendAnchor(plugin, pOutgoingConnectionName);
    }
    appendOutgoingConnection(anchor, pIncConnInt);
}

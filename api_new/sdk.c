#include "sdk.h"

const int cacheSize = 4194304; //4mb

/*
** The structure of a plugin handle looks like this:
**
** (struct PluginSharedMemory)
**     toolId (uint32_t)
**     toolConfig (wchar_t *)
**     engine (struct EngineInterface*)
**     outputAnchors (struct OutputAnchor*)
**         name (char *)
**         metadata (wchar_t *)
**         isOpen (uint32_t)
**         firstChild (struct OutputConn*)
**             isOpen (uint32_t)
**             ii (struct IncomingInterface*)
**             nextConnection (struct OutputConn*)
**         nextAnchor (struct OutputAnchor*)
**         recordCache (char *)
**         recordCachePosition (uint32_t)
**     totalInputConnections (uint32_t)
**     closedInputConnections (uint32_t)
**     inputAnchors (struct InputAnchor*)
**         name (wchar_t *)
**         firstChild (struct InputConnection*)
**             isOpen (uint32_t)
**             metadata (wchar_t *)
**             percent (double)
**             nextConnection (struct InputConnection*)
**             plugin (struct PluginSharedMemory*)
**             fixedFieldSize (uint32_t)
**             varFieldSize (uint32_t)
**             recordCache (char *)
**             recordCachePosition (uint32_t)
**         nextAnchor (struct InputAnchor*)
*/

struct InputConnection {
    uint32_t                   isOpen;
    wchar_t*                   metadata;
    double                     percent;
    struct InputConnection*    nextConnection;
    struct PluginSharedMemory* plugin;
    uint32_t                   fixedFieldSize;
    uint32_t                   varFieldSize;
    char*                      recordCache;
    uint32_t                   recordCachePosition;
};

struct InputAnchor {
    wchar_t*                name;
    struct InputConnection* firstChild;
    struct InputAnchor*     nextAnchor;
};

struct OutputConn {
    uint32_t                            isOpen;
    struct IncomingConnectionInterface* ii;
    struct OutputConn*                  nextConnection;
};

struct OutputAnchor {
    wchar_t*             name;
    wchar_t*             metadata;
    uint32_t             isOpen;
    struct OutputConn*   firstChild;
    struct OutputAnchor* nextAnchor;
    char*                recordCache;
    uint32_t             recordCachePosition;
};

struct PluginSharedMemory {
    uint32_t                toolId;
    wchar_t*                toolConfig;
    struct EngineInterface* engine;
    struct OutputAnchor*    outputAnchors;
    uint32_t                totalInputConnections;
    uint32_t                closedInputConnections;
    struct InputAnchor*     inputAnchors;
};

void* configurePlugin(uint32_t nToolID, wchar_t * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface) {
    struct PluginSharedMemory* plugin = malloc(sizeof(struct PluginSharedMemory));
    plugin->toolId = nToolID;
    plugin->toolConfig = pXmlProperties;
    plugin->engine = pEngineInterface;
    plugin->outputAnchors = NULL;
    plugin->totalInputConnections = 0;
    plugin->closedInputConnections = 0;
    plugin->inputAnchors = NULL;

    r_pluginInterface->handle = plugin;
    r_pluginInterface->pPI_Close = &PI_Close;
    r_pluginInterface->pPI_PushAllRecords = &PI_PushAllRecords;
    r_pluginInterface->pPI_AddIncomingConnection = &PI_AddIncomingConnection;
    r_pluginInterface->pPI_AddOutgoingConnection = &PI_AddOutgoingConnection;

    Init(plugin);

    return plugin;
}

void PI_Close(void * handle, bool bHasErrors) {
    // do nothing
}

long PI_PushAllRecords(void * handle, __int64 nRecordLimit){
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    OnComplete(plugin);
}

struct InputAnchor* createInputAnchor(wchar_t* name) {
    struct InputAnchor* anchor = malloc(sizeof(struct InputAnchor));
    anchor->name = name;
    anchor->firstChild = NULL;
    anchor->nextAnchor = NULL;
    return anchor;
}

struct InputAnchor* getOrCreateInputAnchor(struct PluginSharedMemory* plugin, wchar_t* name) {
    if (NULL == plugin->inputAnchors) {
        struct InputAnchor* anchor = createInputAnchor(name);
        plugin->inputAnchors = anchor;
        return anchor;
    }

    struct InputAnchor* anchor = plugin->inputAnchors;
    while (true) {
        if (wcscmp(name, anchor->name) == 0) {
            return anchor;
        }
        if (NULL == anchor->nextAnchor) {
            break;
        }
        anchor = anchor->nextAnchor;
    }

    struct InputAnchor* child = createInputAnchor(name);
    anchor->nextAnchor = child;
    return child;
}

long PI_AddIncomingConnection(void * handle, wchar_t * pIncomingConnectionType, wchar_t * pIncomingConnectionName, struct IncomingConnectionInterface *r_IncConnInt) {
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    struct InputAnchor *anchor = getOrCreateInputAnchor(plugin, pIncomingConnectionName);
    struct InputConnection *connection = malloc(sizeof(struct InputConnection));
    connection->isOpen = 1;
    connection->metadata = NULL;
    connection->percent = 0;
    connection->nextConnection = NULL;
    connection->plugin = plugin;
    connection->fixedFieldSize = 0;
    connection->varFieldSize = 0;
    connection->recordCache = NULL;
    connection->recordCachePosition = 0;

    plugin->totalInputConnections++;

    r_IncConnInt->handle = connection;
    r_IncConnInt->pII_Init = &II_Init;
    r_IncConnInt->pII_PushRecord = &II_PushRecord;
    r_IncConnInt->pII_UpdateProgress = &II_UpdateProgress;
    r_IncConnInt->pII_Close = &II_Close;
    r_IncConnInt->pII_Free = &II_Free;

    return 1;
}

struct OutputAnchor* getOutputAnchorByName(struct OutputAnchor* anchor, wchar_t* name) {
    if (NULL == anchor) {
        return NULL;
    }
    if (wcscmp(name, anchor->name) == 0) {
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

struct OutputAnchor* appendOutgoingAnchor(struct PluginSharedMemory* plugin, wchar_t * name) {
    struct OutputAnchor* anchor = malloc(sizeof(struct OutputAnchor));
    anchor->name = name;
    anchor->metadata = NULL;
    anchor->isOpen = 0;
    anchor->firstChild = NULL;
    anchor->nextAnchor = NULL;
    anchor->recordCache = NULL;
    anchor->recordCachePosition = 0;

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

long PI_AddOutgoingConnection(void * handle, wchar_t * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt) {
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    struct OutputAnchor* anchor = getOutputAnchorByName(plugin->outputAnchors, pOutgoingConnectionName);
    if (NULL == anchor) {
        anchor = appendOutgoingAnchor(plugin, pOutgoingConnectionName);
    }
    appendOutgoingConnection(anchor, pIncConnInt);
}

long II_Init(void * handle, wchar_t * pXmlRecordMetaInfo) {
    struct InputConnection *input = (struct InputConnection*)handle;
    input->metadata = pXmlRecordMetaInfo;
    OnInputConnectionOpened(input);
    return 1;
}

uint32_t uint32FromRecordPosition(char * record, uint32_t position) {
    uint32_t* value = (uint32_t*)(&record[position]);
    return *value;
}

long II_PushRecord(void * handle, char * pRecord) {
    struct InputConnection *input = (struct InputConnection*)handle;
    uint32_t totalSize = input->fixedFieldSize + input->varFieldSize;
    if (input->varFieldSize > 0) {
        uint32_t varSize = uint32FromRecordPosition(pRecord, totalSize);
        totalSize += varSize;
    }

    if (input->recordCachePosition + totalSize > cacheSize && input->recordCachePosition > 0) {
        OnRecordPacket(input);
        input->recordCachePosition = 0;
    }

    if (totalSize > cacheSize) {
        OnSingleRecord(input, pRecord);
        return 1;
    }

    memcpy(&input->recordCache[input->recordCachePosition], pRecord, totalSize);
    input->recordCachePosition += totalSize;
    return 1;
}

void II_UpdateProgress(void * handle, double dPercent) {
    struct InputConnection *input = (struct InputConnection*)handle;
    input->percent = dPercent;
}

void II_Close(void * handle) {
    struct InputConnection *input = (struct InputConnection*)handle;
    struct PluginSharedMemory *plugin = input->plugin;
    plugin->closedInputConnections++;

    if (plugin->totalInputConnections != plugin->closedInputConnections) {
        return;
    }
    OnComplete(plugin);
}

void II_Free(void * handle) {

}

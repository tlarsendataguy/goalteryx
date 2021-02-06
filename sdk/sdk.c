#include "sdk.h"

const int cacheSize = 4194304; //4mb
const int STATUS_Complete = 4;
const int STATUS_UpdateOutputMetaInfoXml = 10;
const int STATUS_RecordCountString = 50;
utf16char empty[1] = {0};

/*
** The structure of a plugin handle looks like this:
**
** (struct PluginSharedMemory)
**     toolId (uint32_t)
**     toolConfig (utf16char *)
**     toolConfigLen (uint32_t)
**     engine (struct EngineInterface*)
**     ayxInterface (struct PluginInterface*)
**     outputAnchors (struct OutputAnchor*)
**         name (utf16char *)
**         metadata (utf16char *)
**         browseEverywhereId (uint32_t)
**         isOpen (char)
**         plugin (struct PluginSharedMemory*)
**         firstChild (struct OutputConn*)
**             isOpen (char)
**             ii (struct IncomingInterface*)
**             nextConnection (struct OutputConn*)
**         nextAnchor (struct OutputAnchor*)
**         fixedSize (uint32_t)
**         hasVarFields (char)
**         recordCache (char *)
**         recordCachePosition (uint32_t)
**         recordCacheSize (uint32_t)
**         recordCount (uint64_t)
**         totalDataSize (uint64_t)
**     totalInputConnections (uint32_t)
**     closedInputConnections (uint32_t)
**     inputAnchors (struct InputAnchor*)
**         name (utf16char *)
**         firstChild (struct InputConnection*)
**             anchor (struct InputAnchor*)
**             isOpen (char)
**             metadata (utf16char *)
**             percent (double)
**             nextConnection (struct InputConnection*)
**             plugin (struct PluginSharedMemory*)
**             fixedSize (uint32_t)
**             hasVarFields (char)
**             recordCache (char *)
**             recordCachePosition (uint32_t)
**             recordCacheSize (uint32_t)
**         nextAnchor (struct InputAnchor*)
*/

struct PluginInterface* generatePluginInterface(){
    return malloc(sizeof(struct PluginInterface));
}

struct IncomingConnectionInterface* generateIncomingConnectionInterface(){
    return malloc(sizeof(struct IncomingConnectionInterface));
}

void callPiAddIncomingConnection(struct PluginSharedMemory *handle, utf16char * name, struct IncomingConnectionInterface *ii){
    PI_AddIncomingConnection(handle, empty, name, ii);
}

void callPiAddOutgoingConnection(struct PluginSharedMemory *handle, utf16char * name, struct IncomingConnectionInterface *ii){
    PI_AddOutgoingConnection(handle, name, ii);
}

void simulateInputLifecycle(struct PluginInterface *pluginInterface) {
    pluginInterface->pPI_PushAllRecords(pluginInterface->handle, 0);
    pluginInterface->pPI_Close(pluginInterface->handle, 0);
}

void sendMessage(struct EngineInterface * engine, int nToolID, int nStatus, utf16char *pMessage){
    if (NULL != engine) {
        engine->pOutputMessage(engine->handle, nToolID, nStatus, pMessage);
    }
}

void outputToolProgress(struct EngineInterface * engine, int nToolID, double progress){
    engine->pOutputToolProgress(engine->handle, nToolID, progress);
}

void sendProgressToAnchor(struct OutputAnchor *anchor, double progress) {
    struct OutputConn *conn = anchor->firstChild;
    while (conn != NULL) {
        if (conn->isOpen == 1) {
            conn->ii->pII_UpdateProgress(conn->ii->handle, progress);
        }
        conn = conn->nextConnection;
    }
}

void* getInitVar(struct EngineInterface * engine, utf16char *pVar) {
    return engine->pGetInitVar(engine->handle, pVar);
}

uint32_t getLenFromUtf16Ptr(utf16char * ptr) {
    uint32_t len = 0;
    while (ptr[len] != 0) {
        len++;
    }
    return len;
}

void* configurePlugin(uint32_t nToolID, utf16char * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface) {
    struct PluginSharedMemory* plugin = malloc(sizeof(struct PluginSharedMemory));
    plugin->toolId = nToolID;
    plugin->toolConfig = pXmlProperties;
    plugin->toolConfigLen = getLenFromUtf16Ptr(pXmlProperties);
    plugin->engine = pEngineInterface;
    plugin->ayxInterface = r_pluginInterface;
    plugin->outputAnchors = NULL;
    plugin->totalInputConnections = 0;
    plugin->closedInputConnections = 0;
    plugin->inputAnchors = NULL;

    r_pluginInterface->handle = plugin;
    r_pluginInterface->pPI_Close = &PI_Close;
    r_pluginInterface->pPI_PushAllRecords = &PI_PushAllRecords;
    r_pluginInterface->pPI_AddIncomingConnection = &PI_AddIncomingConnection;
    r_pluginInterface->pPI_AddOutgoingConnection = &PI_AddOutgoingConnection;

    return plugin;
}

void appendOutgoingConnection(struct OutputAnchor* anchor, struct IncomingConnectionInterface* ii) {
    struct OutputConn* conn = malloc(sizeof(struct OutputConn));
    conn->isOpen = 0;
    conn->ii = ii;
    conn->nextConnection = NULL;

    if (NULL == anchor->firstChild) {
        anchor->firstChild = conn;
    } else {
        struct OutputConn *childConn = anchor->firstChild;
        while (childConn->nextConnection != NULL) {
            childConn = childConn->nextConnection;
        }
        childConn->nextConnection = conn;
    }

    if (anchor->isOpen == 1) {
        long result = ii->pII_Init(ii->handle, anchor->metadata);
        if (result == 1) {
            conn->isOpen = 1;
        }
    }
    return;
}

void openOutgoingAnchor(struct OutputAnchor *anchor, utf16char * config) {
    anchor->metadata = config;
    struct EngineInterface* engine = anchor->plugin->engine;
    sendMessage(engine, anchor->plugin->toolId, STATUS_UpdateOutputMetaInfoXml, config);

    if (anchor->plugin->engine != NULL && anchor->browseEverywhereId > 0) {
        struct EngineInterface* engine = anchor->plugin->engine;
        struct IncomingConnectionInterface* ii = engine->pBrowseEverywhereGetII(engine->handle, anchor->browseEverywhereId, anchor->plugin->toolId, anchor->name);
        appendOutgoingConnection(anchor, ii);
    }

    anchor->isOpen = 1;
    struct OutputConn * conn = anchor->firstChild;
    while (NULL != conn) {
        long result = conn->ii->pII_Init(conn->ii->handle, config);
        if (result == 1) {
            conn->isOpen = 1;
        }
        conn = conn->nextConnection;
    }
}

void PI_Close(void * handle, bool bHasErrors) {
    // do nothing
}

void closeAllOutputAnchors(struct OutputAnchor *anchor) {
    while (anchor != NULL) {
        struct OutputConn *conn = anchor->firstChild;
        while (anchor->isOpen == 1 && conn != NULL) {
            if (conn->isOpen == 1) {
                conn->ii->pII_Close(conn->ii->handle);
                conn->isOpen = 0;
            }
            conn = conn->nextConnection;
        }
        anchor = anchor->nextAnchor;
    }
}

void freeAllInputAnchors(struct InputAnchor *anchor) {
    struct InputAnchor *nextAnchor;
    struct InputConnection *connection;
    struct InputConnection *nextConnection;

    while (anchor != NULL) {
        nextAnchor = anchor->nextAnchor;
        connection = anchor->firstChild;
        while (connection != NULL) {
            nextConnection = connection->nextConnection;

            free(connection->metadata);
            free(connection);

            connection = nextConnection;
        }

        free(anchor);
        anchor = nextAnchor;
    }
}

void freeAllOutputAnchors(struct OutputAnchor *anchor) {
    struct OutputAnchor *nextAnchor;
    struct OutputConn *connection;
    struct OutputConn *nextConnection;

    while (anchor != NULL) {
        nextAnchor = anchor->nextAnchor;

        connection = anchor->firstChild;
        while (connection != NULL) {
            nextConnection = connection->nextConnection;
            free(connection);
            connection = nextConnection;
        }

        free(anchor->metadata);
        free(anchor->recordCache);
        free(anchor);
        anchor = nextAnchor;
    }
}

void complete(struct PluginSharedMemory *plugin) {
    goOnComplete(plugin);
    freeAllInputAnchors(plugin->inputAnchors);
    closeAllOutputAnchors(plugin->outputAnchors);
    freeAllOutputAnchors(plugin->outputAnchors);
    sendMessage(plugin->engine, plugin->toolId, STATUS_Complete, empty);
    //free(plugin->toolConfig);
    free(plugin);
}

long PI_PushAllRecords(void * handle, int64_t nRecordLimit){
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    complete(plugin);
    return 1;
}

struct InputAnchor* createInputAnchor(utf16char* name) {
    struct InputAnchor* anchor = malloc(sizeof(struct InputAnchor));
    anchor->name = name;
    anchor->firstChild = NULL;
    anchor->nextAnchor = NULL;
    return anchor;
}

bool isUtf16Equal(utf16char* first, utf16char* second) {
    int index = 0;
    while (true) {
        if (first[index] != second[index]) {
            return false;
        }
        if (first[index] == 0) {
            return true;
        }
        index++;
    }
}

struct InputAnchor* getOrCreateInputAnchor(struct PluginSharedMemory* plugin, utf16char* name) {
    if (NULL == plugin->inputAnchors) {
        struct InputAnchor* anchor = createInputAnchor(name);
        plugin->inputAnchors = anchor;
        return anchor;
    }

    struct InputAnchor* anchor = plugin->inputAnchors;
    while (true) {
        if (isUtf16Equal(name, anchor->name)) {
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

long PI_AddIncomingConnection(void * handle, utf16char * pIncomingConnectionType, utf16char * pIncomingConnectionName, struct IncomingConnectionInterface *r_IncConnInt) {
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    struct InputAnchor *anchor = getOrCreateInputAnchor(plugin, pIncomingConnectionName);
    struct InputConnection *connection = malloc(sizeof(struct InputConnection));
    connection->anchor = anchor;
    connection->isOpen = 1;
    connection->metadata = NULL;
    connection->percent = 0;
    connection->nextConnection = NULL;
    connection->plugin = plugin;
    connection->fixedSize = 0;
    connection->hasVarFields = 0;
    connection->recordCache = NULL;
    connection->recordCachePosition = 0;
    connection->recordCacheSize = 0;

    plugin->totalInputConnections++;

    r_IncConnInt->handle = connection;
    r_IncConnInt->pII_Init = &II_Init;
    r_IncConnInt->pII_PushRecord = &II_PushRecord;
    r_IncConnInt->pII_UpdateProgress = &II_UpdateProgress;
    r_IncConnInt->pII_Close = &II_Close;
    r_IncConnInt->pII_Free = &II_Free;

    return 1;
}

struct OutputAnchor* getOutputAnchorByName(struct OutputAnchor* anchor, utf16char* name) {
    if (NULL == anchor) {
        return NULL;
    }
    if (isUtf16Equal(name, anchor->name)) {
        return anchor;
    }
    return getOutputAnchorByName(anchor->nextAnchor, name);
}

struct OutputAnchor* createOutgoingAnchor(utf16char* name) {
    struct OutputAnchor* anchor = malloc(sizeof(struct OutputAnchor));
    anchor->name = name;
    anchor->metadata = NULL;
    anchor->browseEverywhereId = 0;
    anchor->isOpen = 0;
    anchor->firstChild = NULL;
    anchor->nextAnchor = NULL;
    anchor->fixedSize = 0;
    anchor->hasVarFields = 0;
    anchor->recordCache = NULL;
    anchor->recordCachePosition = 0;
    anchor->recordCacheSize = 0;
    anchor->recordCount = 0;
    anchor->totalDataSize = 0;

    return anchor;
}

struct OutputAnchor* appendOutgoingAnchor(struct PluginSharedMemory* plugin, utf16char * name) {
    struct OutputAnchor* anchor = createOutgoingAnchor(name);
    anchor->plugin = plugin;
    if (plugin->engine != NULL) {
        anchor->browseEverywhereId = plugin->engine->pBrowseEverywhereReserveAnchor(plugin->engine->handle, plugin->toolId);
    }

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

long PI_AddOutgoingConnection(void * handle, utf16char * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt) {
    struct PluginSharedMemory *plugin = (struct PluginSharedMemory*)handle;
    struct OutputAnchor* anchor = getOutputAnchorByName(plugin->outputAnchors, pOutgoingConnectionName);
    if (NULL == anchor) {
        anchor = appendOutgoingAnchor(plugin, pOutgoingConnectionName);
    }
    appendOutgoingConnection(anchor, pIncConnInt);
    return 1;
}

long II_Init(void * handle, utf16char * pXmlRecordMetaInfo) {
    struct InputConnection *input = (struct InputConnection*)handle;
    input->metadata = pXmlRecordMetaInfo;
    goOnInputConnectionOpened(input);
    return 1;
}

uint32_t uint32FromRecordPosition(char * record, uint32_t position) {
    uint32_t* value = (uint32_t*)(&record[position]);
    return *value;
}

long II_PushRecord(void * handle, char * pRecord) {
    struct InputConnection *input = (struct InputConnection*)handle;
    uint32_t totalSize = input->fixedSize;
    if (input->hasVarFields == 1) {
        uint32_t varSize = uint32FromRecordPosition(pRecord, totalSize);
        totalSize += 4 + varSize;
    }

    if (totalSize > input->recordCacheSize) {
        if (input->recordCachePosition > 0) {
            goOnRecordPacket(handle);
            input->recordCachePosition = 0;
        }

        if (input->recordCacheSize > 0) {
            free(input->recordCache);
        }

        uint32_t newCacheSize = cacheSize;
        if (totalSize > newCacheSize) {
            newCacheSize = totalSize;
        }

        input->recordCache = malloc(newCacheSize);
        input->recordCacheSize = newCacheSize;
    }

    if (input->recordCachePosition + totalSize > input->recordCacheSize) {
        goOnRecordPacket(handle);
        input->recordCachePosition = 0;
    }

    memcpy(input->recordCache+input->recordCachePosition, pRecord, totalSize);
    input->recordCachePosition += totalSize;
    return 1;
}

void II_UpdateProgress(void * handle, double dPercent) {
    struct InputConnection *input = (struct InputConnection*)handle;
    input->percent = dPercent;
}

void II_Close(void * handle) {
    struct InputConnection *input = (struct InputConnection*)handle;
    if (input->recordCachePosition > 0) {
        goOnRecordPacket(input);
    }
    struct PluginSharedMemory *plugin = input->plugin;
    plugin->closedInputConnections++;

    free(input->recordCache);

    if (plugin->totalInputConnections != plugin->closedInputConnections) {
        return;
    }
    complete(plugin);
}

void II_Free(void * handle) {

}

const utf16char pipe = 124;
const utf16char zero = 48;

int uint64ToString(utf16char cache[20], uint64_t value) { // the largest uint64 value is 20 digits long, so utf16char[20] is large enough for our purposes
    int index = 0;
    if (value == 0) {
        cache[0] = zero;
        return 1;
    }
    while (value != 0) {
        utf16char rem = value % 10;
        cache[index] = zero + rem;
        index++;
        value = value/10;
    }
    return index;
}

void formatRecordCountString(utf16char* cache, size_t len, utf16char* outputName, uint64_t recordCount, uint64_t totalDataSize) {
    int index = 0;

    // copy output name
    int nameIndex = 0;
    while (outputName[nameIndex] != 0) {
        if (index >= len-1) {
            break;
        }
        cache[index] = outputName[nameIndex];
        index++;
        nameIndex++;
    }

    cache[index] = pipe;
    index++;

    // copy record count
    utf16char integerCache[20];
    int endIndex = uint64ToString(integerCache, recordCount) - 1;
    while (endIndex >= 0) {
        if (index >= len-1) {
            break;
        }
        cache[index] = integerCache[endIndex];
        index++;
        endIndex--;
    }

    cache[index] = pipe;
    index++;

    // copy data size
    endIndex = uint64ToString(integerCache, totalDataSize) - 1;
    while (endIndex >= 0) {
        if (index >= len-1) {
            break;
        }
        cache[index] = integerCache[endIndex];
        index++;
        endIndex--;
    }

    cache[index] = 0; // null terminator
}

void callWriteRecords(struct OutputAnchor *anchor) {
    struct OutputConn *conn = anchor->firstChild;
    if (NULL == conn) {
        anchor->recordCachePosition = 0;
        return;
    }
    char *record;
    uint32_t written = 0;
    while (written < anchor->recordCachePosition) {
        conn = anchor->firstChild;
        record = &anchor->recordCache[written];
        while (conn != NULL) {
            if (conn->isOpen == 0) {
                conn = conn->nextConnection;
                continue;
            }
            long result = conn->ii->pII_PushRecord(conn->ii->handle, record);
            if (result == 0) {
                conn->ii->pII_Close(conn->ii->handle);
                conn->isOpen = 0;
            }
            conn = conn->nextConnection;
        }
        written += anchor->fixedSize;
        if (anchor->hasVarFields == 1) {
            uint32_t varLen = uint32FromRecordPosition(anchor->recordCache, written);
            written += 4 + varLen;
        }
        anchor->recordCount++;
    }
    anchor->totalDataSize += written;
    utf16char msg[128];
    formatRecordCountString(msg, sizeof(msg), anchor->name, anchor->recordCount, anchor->totalDataSize);
    sendMessage(anchor->plugin->engine, anchor->plugin->toolId, STATUS_RecordCountString, &msg[0]);
}

void* allocateCache(int size) {
    return malloc(size);
}
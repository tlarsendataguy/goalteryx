#include <stdlib.h>
#include "plugins.h"


// Plugin methods

// We can have up to 1,000 incoming interfaces.  We store information about those interfaces in these variables.
// buffers contains the 10-record buffer for each incoming interface we add.  iiFixedSizes contains the fixed size
// of each incoming interface.  Both are accessed using the index-based handle defined in go_piAddIncomingConnection.
int currentIiIndex = 0;
struct IncomingRecordCache* buffers[1000];
int iiFixedSizes[1000];

// a reference to the engine provided to our plugins.
struct EngineInterface *engine;

// wires up the PluginInterface struct with the C functions below.  The C functions call into the Go layer.
void c_configurePlugin(void * handle, struct PluginInterface * pluginInterface) {
    pluginInterface->handle = handle;
    pluginInterface->pPI_PushAllRecords = c_piPushAllRecords;
    pluginInterface->pPI_AddIncomingConnection = c_piAddIncomingConnection;
    pluginInterface->pPI_AddOutgoingConnection = c_piAddOutgoingConnection;
    pluginInterface->pPI_Close = c_piClose;
}

// C entry point for PI_PushAllRecords.  Calls into Go.
long c_piPushAllRecords(void * handle, __int64 nRecordLimit) {
    return go_piPushAllRecords(handle, nRecordLimit);
}

// C entry point for PI_Close.  Calls into Go.
void c_piClose(void * handle, bool bHasErrors) {
    go_piClose(handle, bHasErrors);
}

// C entry point for PI_AddIncomingConnection.  Calls into Go.  If Go returns a PreSort XML string, call PreSort on
// the engine.  Wire up the appropriate IncomingInterface struct with the relevant C II functions.
long c_piAddIncomingConnection(void * handle, void * connectionType, void * connectionName, struct IncomingConnectionInterface * incomingInterface) {
    struct IncomingConnectionInfo *info = go_piAddIncomingConnection(handle, connectionType, connectionName);
    if (!info) {
        return 0;
    }

    struct IncomingConnectionInterface *actualIncomingInterface;
    if (info->presortString) {
        struct PreSortConnectionInterface *newPresortTool;
        long result = engine->pPreSort(engine->handle, 1, info->presortString, incomingInterface, &actualIncomingInterface, &newPresortTool);
    } else {
        actualIncomingInterface = incomingInterface;
    }

    actualIncomingInterface->handle = info->handle;
    actualIncomingInterface->pII_Init = c_iiInit;

    int cacheSize = info->cacheSize;
    if (cacheSize == 0) {
        actualIncomingInterface->pII_PushRecord = c_iiPushRecord;
    } else {
        actualIncomingInterface->pII_PushRecord = c_iiPushRecordCache;
        struct IncomingRecordCache* cache = malloc(sizeof(struct IncomingRecordCache));
        cache->recordsInBuffer = cacheSize;
        cache->buffer = malloc(sizeof(void*)*cacheSize);
        cache->bufferSizes = malloc(sizeof(int)*cacheSize);

        int iiIndex = *((int*)info->handle);
        buffers[iiIndex] = cache;
    }

    actualIncomingInterface->pII_UpdateProgress = c_iiUpdateProgress;
    actualIncomingInterface->pII_Close = c_iiClose;
    actualIncomingInterface->pII_Free = c_iiFree;
    free(info);
    return 1;
}

// C entry point for PI_AddOutgoingConnection.  Calls into Go.
long c_piAddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt){
    return go_piAddOutgoingConnection(handle, pOutgoingConnectionName, pIncConnInt);
}

// Create the IncomingConnectionInfo struct which contains an IncomingInterface handle and a PreSort string.  Go calls
// this function from go_piAddIncomingConnection and populates this struct.  This struct is then returned back to
// c_piAddIncomingConnection.
struct IncomingConnectionInfo *newSortedIncomingConnectionInfo(void * handle, void * presortString, int cacheSize){
    struct IncomingConnectionInfo *info = malloc(sizeof(struct IncomingConnectionInfo));
    info->handle = handle;
    info->presortString = presortString;
    return info;
}

// Create the IncomingConnectionInfo struct which contains an IncomingInterface handle with no PreSort.  Go calls
// this function from go_piAddIncomingConnection and populates this struct.  This struct is then returned back to
// c_piAddIncomingConnection.
struct IncomingConnectionInfo *newUnsortedIncomingConnectionInfo(void * handle, int cacheSize){
    struct IncomingConnectionInfo *info = malloc(sizeof(struct IncomingConnectionInfo));
    info->handle = handle;
    return info;
}


// Incoming interface methods

// Used to generate an IncomingInterface struct for testing outside of Alteryx.
struct IncomingConnectionInterface* newIi(void * iiHandle) {
    struct IncomingConnectionInterface *ptr;
    ptr = malloc(sizeof(struct IncomingConnectionInterface));
    ptr->handle = iiHandle;
    ptr->pII_Init = c_iiInit;
    ptr->pII_PushRecord = c_iiPushRecord;
    ptr->pII_UpdateProgress = c_iiUpdateProgress;
    ptr->pII_Close = c_iiClose;
    ptr->pII_Free = c_iiFree;
    return ptr;
}

// Generate a pointer to our incoming interface index handle.
void * getIiIndex(){
    int* iiIndex = malloc(sizeof(int));
    *iiIndex = currentIiIndex++;
    return iiIndex;
}

void saveIncomingInterfaceFixedSize(void * handle, int fixedSize) {
    int iiIndex = *((int*)handle);
    iiFixedSizes[iiIndex] = fixedSize;
}

// C entry point for II_Init.  Calls into Go.
long c_iiInit(void * handle, void * recordInfoIn) {
    return go_iiInit(handle, recordInfoIn);
}

// C entry point for II_PushRecord.  This function buffers records and only calls into Go when the buffer has filled.
long c_iiPushRecordCache(void * handle, void * record) {
    int iiIndex = *((int*)handle);
    int fixedSize = iiFixedSizes[iiIndex];
    struct IncomingRecordCache *buffer = buffers[iiIndex];
    if (buffer->currentBufferIndex == buffer->recordsInBuffer) {
        go_iiPushRecordCache(handle, buffer->buffer, buffer->currentBufferIndex);
        buffer->currentBufferIndex = 0;
    }

    int varSize = *(int*)(record+fixedSize);
    int totalSize = fixedSize + 4 + varSize;
    int bufferSize = (*buffer->bufferSizes)[buffer->currentBufferIndex];
    if (totalSize > bufferSize) {
        if (bufferSize > 0) {
            free((*buffer->buffer)[buffer->currentBufferIndex]);
        }
        (*buffer->buffer)[buffer->currentBufferIndex] = malloc(totalSize);
        (*buffer->bufferSizes)[buffer->currentBufferIndex] = totalSize;
    }
    memcpy((*buffer->buffer)[buffer->currentBufferIndex], record, totalSize);
    buffer->currentBufferIndex++;
    buffer->recordCount++;
    return 1;
}

// C entry point for directly pushing records to Go, without a buffer.  Certain classes of tools must receive records
// immediately rather than wait for the cache to fill.
long c_iiPushRecord(void * handle, void * record) {
    return go_iiPushRecord(handle, record);
}

// C entry point for II_UpdateProgress.  Calls into Go.
void c_iiUpdateProgress(void * handle, double percent){
    go_iiUpdateProgress(handle, percent);
}

// C entry point for II_Close.  Calls into Go.
void c_iiClose(void * handle){
    int iiIndex = *((int*)handle);
    struct IncomingRecordCache *buffer = buffers[iiIndex];
    if (buffer) {
        go_iiPushRecordCache(handle, buffer->buffer, buffer->currentBufferIndex);
        buffer->currentBufferIndex = 0;
    }
    go_iiClose(handle);
}

// C entry point for II_Free.  This frees the buffer.  We do not call into Go; Go should cleanup when Close is called.
void c_iiFree(void * handle){
    int iiIndex = *((int*)handle);
    struct IncomingRecordCache *buffer = buffers[iiIndex];

    if (buffer){
        int ceiling = buffer->recordsInBuffer;
        if (buffer->recordCount < ceiling) {
            ceiling = buffer->recordCount;
        }

        for (int i = 0; i < ceiling; i++) {
            free((*buffer->buffer)[i]);
        }

        free(buffer->buffer);
        free(buffer->bufferSizes);
        free(buffer);
    }
    free(handle);
}


// Output connection methods

long c_outputInit(struct IncomingConnectionInterface * connection, void * recordMetaInfoXml) {
    return connection->pII_Init(connection->handle, recordMetaInfoXml);
}

long c_outputPushRecord(struct IncomingConnectionInterface * connection, void * record) {
    return connection->pII_PushRecord(connection->handle, record);
}

long c_outputClose(struct IncomingConnectionInterface * connection) {
    connection->pII_Close(connection->handle);
}

void c_outputUpdateProgress(struct IncomingConnectionInterface * connection, double percent){
    connection->pII_UpdateProgress(connection->handle, percent);
}


// Engine methods

void c_setEngine(struct EngineInterface *pEngineInterface) {
    engine = pEngineInterface;
}

void callEngineOutputMessage(int toolId, int status, void * message) {
	engine->pOutputMessage(engine->handle, toolId, status, message);
}

void * callEngineCreateTempFileName(void * ext) {
	return engine->pCreateTempFileName(engine->handle, ext);
}

unsigned callEngineBrowseEverywhereReserveAnchor(int toolId) {
	return engine->pBrowseEverywhereReserveAnchor(engine->handle, toolId);
}

long callEngineOutputToolProgress(int toolId, double dPercentProgress) {
    return engine->pOutputToolProgress(engine->handle, toolId, dPercentProgress);
}

struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(unsigned browseEverywhereAnchorId, int toolId, void * name) {
	return engine->pBrowseEverywhereGetII(engine->handle, browseEverywhereAnchorId, toolId, name);
}

void * callEngineGetInitVar(void * initVar) {
    void *value = engine->pGetInitVar(engine->handle, initVar);
    return value;
}

void * callEngineGetInitVar2(int toolId, void * initVar) {
    void *value = engine->pGetInitVar2(engine->handle, toolId, initVar);
    return value;
}

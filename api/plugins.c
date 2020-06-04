#include <stdlib.h>
#include <stdio.h>
#include "plugins.h"


// Plugin methods

struct EngineInterface *engine;

void c_configurePlugin(void * handle, struct PluginInterface * pluginInterface) {
    pluginInterface->handle = handle;
    pluginInterface->pPI_PushAllRecords = c_piPushAllRecords;
    pluginInterface->pPI_AddIncomingConnection = c_piAddIncomingConnection;
    pluginInterface->pPI_AddOutgoingConnection = c_piAddOutgoingConnection;
    pluginInterface->pPI_Close = c_piClose;
}

long c_piPushAllRecords(void * handle, __int64 nRecordLimit) {
    return go_piPushAllRecords(handle, nRecordLimit);
}

void c_piClose(void * handle, bool bHasErrors) {
    go_piClose(handle, bHasErrors);
}

long c_piAddIncomingConnection(void * handle, void * connectionType, void * connectionName, struct IncomingConnectionInterface * incomingInterface) {
    void * iiHandle = go_piAddIncomingConnection(handle, connectionType, connectionName);

    struct IncomingConnectionInterface *newIncomingInterface;
    struct PreSortConnectionInterface *newPresortTool;

    long result = engine->pPreSort(engine->handle, 1, L"<SortInfo>\n<Field field=\"RowCount\" order=\"Desc\" />\n</SortInfo>\n", incomingInterface, &newIncomingInterface, &newPresortTool);


    newIncomingInterface->handle = iiHandle;
    newIncomingInterface->pII_Init = c_iiInit;
    newIncomingInterface->pII_PushRecord = c_iiPushRecord;
    newIncomingInterface->pII_UpdateProgress = c_iiUpdateProgress;
    newIncomingInterface->pII_Close = c_iiClose;
    newIncomingInterface->pII_Free = c_iiFree;
    return 1;
}

long c_piAddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt){
    return go_piAddOutgoingConnection(handle, pOutgoingConnectionName, pIncConnInt);
}


// Incoming interface methods

int currentIiIndex = 0;
struct IncomingRecordCache* buffers[1000];
int iiFixedSizes[1000];

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

void * getIiIndex(){
    int* iiIndex = malloc(sizeof(int));
    *iiIndex = currentIiIndex++;
    return iiIndex;
}

void saveIncomingInterfaceFixedSize(void * handle, int fixedSize) {
    int iiIndex = *((int*)handle);
    struct IncomingRecordCache* cache = malloc(sizeof(struct IncomingRecordCache));
    buffers[iiIndex] = cache;
    iiFixedSizes[iiIndex] = fixedSize;
}

long c_iiInit(void * handle, void * recordInfoIn) {
    return go_iiInit(handle, recordInfoIn);
}

long c_iiPushRecord(void * handle, void * record) {
    int iiIndex = *((int*)handle);
    int fixedSize = iiFixedSizes[iiIndex];
    struct IncomingRecordCache *buffer = buffers[iiIndex];
    if (buffer->currentBufferIndex == 10) {
        go_iiPushRecordCache(handle, &(buffer->buffer), buffer->currentBufferIndex);
        buffer->currentBufferIndex = 0;
    }

    int varSize = *(int*)(record+fixedSize);
    int totalSize = fixedSize + 4 + varSize;
    int bufferSize = buffer->bufferSizes[buffer->currentBufferIndex];
    if (totalSize > bufferSize) {
        if (bufferSize > 0) {
            free(buffer->buffer[buffer->currentBufferIndex]);
        }
        buffer->buffer[buffer->currentBufferIndex] = malloc(totalSize);
        buffer->bufferSizes[buffer->currentBufferIndex] = totalSize;
    }
    memcpy(buffer->buffer[buffer->currentBufferIndex], record, totalSize);
    buffer->currentBufferIndex++;
    buffer->recordCount++;
    return 1;
}

void c_iiUpdateProgress(void * handle, double percent){
    go_iiUpdateProgress(handle, percent);
}

void c_iiClose(void * handle){
    int iiIndex = *((int*)handle);
    struct IncomingRecordCache *buffer = buffers[iiIndex];
    go_iiPushRecordCache(handle, &(buffer->buffer), buffer->currentBufferIndex);
    buffer->currentBufferIndex = 0;
    go_iiClose(handle);
}

void c_iiFree(void * handle){
    int iiIndex = *((int*)handle);
    struct IncomingRecordCache *buffer = buffers[iiIndex];

    int ceiling = 10;
    if (buffer->recordCount < 10) {
        ceiling = buffer->recordCount;
    }

    for (int i = 0; i < ceiling; i++) {
        free(buffer->buffer[i]);
    }

    free(buffer);
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

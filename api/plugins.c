#include <stdlib.h>
#include <stdio.h>
#include "plugins.h"


// Plugin methods

struct EngineInterface *engine;

void c_configurePlugin(void * handle, struct PluginInterface * pluginInterface, struct EngineInterface * pluginEngine) {
    engine = pluginEngine;
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
    FILE *f = fopen("C:\\temp\\output.txt", "a");
    fprintf(f, "Started c_piAddIncomingConnection\n");
    fflush(f);
    void * iiHandle = go_piAddIncomingConnection(handle, connectionType, connectionName);
    fprintf(f, "Got handle from piAddIncomingConnection\n");
    fflush(f);

    /*
    struct IncomingConnectionInterface *newIncomingInterface;
    struct PreSortConnectionInterface *newPresortTool;

    fprintf(f, "defined variables needed for presort\n");
    fflush(f);
    long result = engine->pPreSort(engine->handle, 1, L"<SortInfo>\n<Field field=\"RowCount\" order=\"Asc\" />\n</SortInfo>\n", incomingInterface, &newIncomingInterface, &newPresortTool);
    fprintf(f, "Got presort\n");
    fflush(f);
    */

    incomingInterface->handle = iiHandle;
    incomingInterface->pII_Init = c_iiInit;
    incomingInterface->pII_PushRecord = c_iiPushRecord;
    incomingInterface->pII_UpdateProgress = c_iiUpdateProgress;
    incomingInterface->pII_Close = c_iiClose;
    incomingInterface->pII_Free = c_iiFree;
    fprintf(f, "Done with c_piAddIncomingConnection\n");
    fclose(f);

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

void callEngineOutputMessage(struct EngineInterface *pEngineInterface, int toolId, int status, void * message) {
	pEngineInterface->pOutputMessage(pEngineInterface->handle, toolId, status, message);
}

void * callEngineCreateTempFileName(struct EngineInterface *pEngineInterface, void * ext) {
	return pEngineInterface->pCreateTempFileName(pEngineInterface->handle, ext);
}

unsigned callEngineBrowseEverywhereReserveAnchor(struct EngineInterface *pEngineInterface, int toolId) {
	return pEngineInterface->pBrowseEverywhereReserveAnchor(pEngineInterface->handle, toolId);
}

long callEngineOutputToolProgress(struct EngineInterface *pEngineInterface, int toolId, double dPercentProgress) {
    return pEngineInterface->pOutputToolProgress(pEngineInterface->handle, toolId, dPercentProgress);
}

struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(struct EngineInterface *pEngineInterface, unsigned browseEverywhereAnchorId, int toolId, void * name) {
	return pEngineInterface->pBrowseEverywhereGetII(pEngineInterface->handle, browseEverywhereAnchorId, toolId, name);
}

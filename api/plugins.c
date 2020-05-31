#include <stdlib.h>
#include "plugins.h"

void callEngineOutputMessage(struct EngineInterface *pEngineInterface, int toolId, int status, void * message) {
	pEngineInterface->pOutputMessage(pEngineInterface->handle, toolId, status, message);
}

void * callEngineCreateTempFileName(struct EngineInterface *pEngineInterface, void * ext) {
	return pEngineInterface->pCreateTempFileName(pEngineInterface->handle, ext);
}

unsigned callEngineBrowseEverywhereReserveAnchor(struct EngineInterface *pEngineInterface, int toolId) {
	return pEngineInterface->pBrowseEverywhereReserveAnchor(pEngineInterface->handle, toolId);
}

struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(struct EngineInterface *pEngineInterface, unsigned browseEverywhereAnchorId, int toolId, void * name) {
	return pEngineInterface->pBrowseEverywhereGetII(pEngineInterface->handle, browseEverywhereAnchorId, toolId, name);
}

long callInitOutput(struct IncomingConnectionInterface * connection, void * recordMetaInfoXml) {
    return connection->pII_Init(connection->handle, recordMetaInfoXml);
}

long callPushRecord(struct IncomingConnectionInterface * connection, void * record) {
    return connection->pII_PushRecord(connection->handle, record);
}

long callCloseOutput(struct IncomingConnectionInterface * connection) {
    connection->pII_Close(connection->handle);
}

struct IncomingConnectionInterface* newIi() {
    struct IncomingConnectionInterface *ptr;
    ptr = malloc(sizeof(struct IncomingConnectionInterface));
    return ptr;
}

void* buffer[10];
int bufferSizes[10] = {0,0,0,0,0,0,0,0,0,0};
int currentBufferIndex = 0;
int currentIiIndex = 0;
int iiFixedSizes[100];

void * getIiIndex(){
    int* iiIndex = malloc(sizeof(int));
    *iiIndex = currentIiIndex++;
    return iiIndex;
}

void saveIncomingInterfaceFixedSize(void * handle, int fixedSize) {
    int iiIndex = *((int*)handle);
    iiFixedSizes[iiIndex] = fixedSize;
}

long iiPushRecord(void * handle, void * record) {
    int iiIndex = *((int*)handle);
    int fixedSize = iiFixedSizes[iiIndex];
    if (currentBufferIndex == 10) {
        pushRecordCache(handle, &buffer, currentBufferIndex);
        currentBufferIndex = 0;
    }

    int varSize = *(int*)(record+fixedSize);
    int totalSize = fixedSize + 4 + varSize;
    int bufferSize = bufferSizes[currentBufferIndex];
    if (totalSize > bufferSize) {
        if (bufferSize > 0) {
            free(buffer[currentBufferIndex]);
        }
        buffer[currentBufferIndex] = malloc(totalSize);
        bufferSizes[currentBufferIndex] = totalSize;
    }
    memcpy(buffer[currentBufferIndex], record, totalSize);
    currentBufferIndex++;
    return 1;
}

void closeRecordCache(void * handle){
    pushRecordCache(handle, &buffer, currentBufferIndex);
    currentBufferIndex = 0;
}

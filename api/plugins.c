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
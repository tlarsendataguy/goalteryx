#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <windows.h>
#include "alteryx_api.h"

void sendMessage(struct EngineInterface * engine, int nToolID, int nStatus, wchar_t *pMessage);
void outputToolProgress(struct EngineInterface * engine, int nToolID, double progress);
void* configurePlugin(uint32_t nToolID, wchar_t * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface);
void PI_Close(void * handle, bool bHasErrors);
long PI_PushAllRecords(void * handle, __int64 nRecordLimit);
long PI_AddIncomingConnection(void * handle,
    wchar_t * pIncomingConnectionType,
    wchar_t * pIncomingConnectionName,
    struct IncomingConnectionInterface *r_IncConnInt);
long PI_AddOutgoingConnection(void * handle,
    wchar_t * pOutgoingConnectionName,
    struct IncomingConnectionInterface *pIncConnInt);
long II_Init(void * handle, wchar_t * pXmlRecordMetaInfo);
long II_PushRecord(void * handle, char * pRecord);
void II_UpdateProgress(void * handle, double dPercent);
void II_Close(void * handle);
void II_Free(void * handle);
void goOnInputConnectionOpened(void * handle);
void goOnRecordPacket(void * handle);
void goOnSingleRecord(void * handle, void * record);
void goOnComplete(void * handle);

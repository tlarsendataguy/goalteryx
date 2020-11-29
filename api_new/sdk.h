#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <windows.h>
#include "alteryx_api.h"

void* configurePlugin(uint32_t nToolID, void * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface);
void PI_Close(void * handle, bool bHasErrors);
long PI_PushAllRecords(void * handle, __int64 nRecordLimit);
long PI_AddIncomingConnection(void * handle,
    void * pIncomingConnectionType,
    void * pIncomingConnectionName,
    struct IncomingConnectionInterface *r_IncConnInt);
long PI_AddOutgoingConnection(void * handle,
    void * pOutgoingConnectionName,
    struct IncomingConnectionInterface *pIncConnInt);
long II_Init(void * handle, void * pXmlRecordMetaInfo);
long II_PushRecord(void * handle, void * pRecord);
void II_UpdateProgress(void * handle, double dPercent);
void II_Close(void * handle);
void II_Free(void * handle);
void Init(void * handle);
void OnInputConnectionOpened(void * handle);
void OnRecordPacket(void * handle);
void OnComplete(void * handle);

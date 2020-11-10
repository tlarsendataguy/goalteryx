#include "alteryx_api.h"

const int cacheSize = 4194304; //4mb

void PI_Close(void * handle, bool bHasErrors);
long PI_PushAllRecords(void * handle, __int64 nRecordLimit);
long PI_AddIncomingConnection(void * handle,
    void * pIncomingConnectionType,
    void * pIncomingConnectionName,
    struct IncomingConnectionInterface *r_IncConnInt);
long PI_AddOutgoingConnection(void * handle,
    void * pOutgoingConnectionName,
    struct IncomingConnectionInterface *pIncConnInt);
long II_Init(void * handle, volatile void * pXmlRecordMetaInfo);
long II_PushRecord(void * handle, volatile void * pRecord);
void II_UpdateProgress(void * handle, double dPercent);
void II_Close(void * handle);
void II_Free(void * handle);

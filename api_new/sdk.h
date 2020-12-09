#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <windows.h>
#include "alteryx_api.h"

struct InputConnection {
    char                       isOpen;
    wchar_t*                   metadata;
    double                     percent;
    struct InputConnection*    nextConnection;
    struct PluginSharedMemory* plugin;
    uint32_t                   fixedSize;
    char                       hasVarFields;
    char*                      recordCache;
    uint32_t                   recordCachePosition;
};

struct InputAnchor {
    wchar_t*                name;
    struct InputConnection* firstChild;
    struct InputAnchor*     nextAnchor;
};

struct OutputConn {
    char                                isOpen;
    struct IncomingConnectionInterface* ii;
    struct OutputConn*                  nextConnection;
};

struct OutputAnchor {
    wchar_t*             name;
    wchar_t*             metadata;
    char                 isOpen;
    struct OutputConn*   firstChild;
    struct OutputAnchor* nextAnchor;
    char*                recordCache;
    uint32_t             recordCachePosition;
};

struct PluginSharedMemory {
    uint32_t                toolId;
    wchar_t*                toolConfig;
    uint32_t                toolConfigLen;
    struct EngineInterface* engine;
    struct OutputAnchor*    outputAnchors;
    uint32_t                totalInputConnections;
    uint32_t                closedInputConnections;
    struct InputAnchor*     inputAnchors;
};

struct PluginInterface* generatePluginInterface();
void simulateInputLifecycle(struct PluginInterface *pluginInterface);
void sendMessage(struct EngineInterface * engine, int nToolID, int nStatus, wchar_t *pMessage);
void outputToolProgress(struct EngineInterface * engine, int nToolID, double progress);
void* getInitVar(struct EngineInterface * engine, wchar_t *pVar);
void* configurePlugin(uint32_t nToolID, wchar_t * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface);
struct OutputAnchor* appendOutgoingAnchor(struct PluginSharedMemory* plugin, wchar_t * name);
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

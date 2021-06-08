#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include "alteryx_api.h"

struct InputConnection {
    struct InputAnchor*        anchor;
    char                       isOpen;
    char                       status;
    utf16char*                 metadata;
    double                     percent;
    struct InputConnection*    nextConnection;
    struct PluginSharedMemory* plugin;
    uint32_t                   fixedSize;
    char                       hasVarFields;
    char*                      recordCache;
    uint32_t                   recordCachePosition;
    uint32_t                   recordCacheSize;
};

struct InputAnchor {
    utf16char*              name;
    struct InputConnection* firstChild;
    struct InputAnchor*     nextAnchor;
};

struct OutputConn {
    char                                isOpen;
    struct IncomingConnectionInterface* ii;
    struct OutputConn*                  nextConnection;
};

struct OutputAnchor {
    utf16char*                 name;
    utf16char*                 metadata;
    uint32_t                   browseEverywhereId;
    char                       isOpen;
    struct PluginSharedMemory* plugin;
    struct OutputConn*         firstChild;
    struct OutputAnchor*       nextAnchor;
    uint32_t                   fixedSize;
    char                       hasVarFields;
    char*                      recordCache;
    uint32_t                   recordCachePosition;
    uint32_t                   recordCacheSize;
    uint64_t                   recordCount;
    uint64_t                   totalDataSize;
};

struct PluginSharedMemory {
    uint32_t                toolId;
    utf16char*              toolConfig;
    uint32_t                toolConfigLen;
    struct EngineInterface* engine;
    struct PluginInterface* ayxInterface;
    struct OutputAnchor*    outputAnchors;
    uint32_t                totalInputConnections;
    uint32_t                closedInputConnections;
    struct InputAnchor*     inputAnchors;
};

struct PluginInterface* generatePluginInterface();
struct IncomingConnectionInterface* generateIncomingConnectionInterface();
void callPiAddIncomingConnection(struct PluginSharedMemory *handle, utf16char * name, struct IncomingConnectionInterface *ii);
void callPiAddOutgoingConnection(struct PluginSharedMemory *handle, utf16char * name, struct IncomingConnectionInterface *ii);
void simulateInputLifecycle(struct PluginInterface *pluginInterface);
void sendMessage(struct EngineInterface * engine, int nToolID, int nStatus, utf16char *pMessage);
void outputToolProgress(struct EngineInterface * engine, int nToolID, double progress);
void sendProgressToAnchor(struct OutputAnchor *anchor, double progress);
void* getInitVar(struct EngineInterface * engine, utf16char *pVar);
void* configurePlugin(uint32_t nToolID, utf16char * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface);
struct OutputAnchor* appendOutgoingAnchor(struct PluginSharedMemory* plugin, utf16char * name);
void openOutgoingAnchor(struct OutputAnchor *anchor, utf16char * config);
void PI_Close(void * handle, bool bHasErrors);
long PI_PushAllRecords(void * handle, int64_t nRecordLimit);
long PI_AddIncomingConnection(void * handle,
    utf16char * pIncomingConnectionType,
    utf16char * pIncomingConnectionName,
    struct IncomingConnectionInterface *r_IncConnInt);
long PI_AddOutgoingConnection(void * handle,
    utf16char * pOutgoingConnectionName,
    struct IncomingConnectionInterface *pIncConnInt);
long II_Init(void * handle, utf16char * pXmlRecordMetaInfo);
long II_PushRecord(void * handle, char * pRecord);
void II_UpdateProgress(void * handle, double dPercent);
void II_Close(void * handle);
void II_Free(void * handle);
void goOnInputConnectionOpened(void * handle);
void goOnRecordPacket(void * handle);
void goOnComplete(void * handle);
void callWriteRecords(struct OutputAnchor *anchor);
void* allocateCache(int size);

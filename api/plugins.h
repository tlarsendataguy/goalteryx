#include <stdbool.h>
#include <stddef.h>
#include <string.h>

// Plugin definitions

struct RecordData
{

};

typedef long ( _stdcall * T_II_Init)(void * handle, void * pXmlRecordMetaInfo);
typedef long ( _stdcall * T_II_PushRecord)(void * handle, void * pRecord);
typedef void ( _stdcall * T_II_UpdateProgress)(void * handle, double dPercent);
typedef void ( _stdcall * T_II_Close)(void * handle);
typedef void ( _stdcall * T_II_Free)(void * handle);

struct IncomingConnectionInterface
{
	int sizeof_IncomingConnectionInterface;
	void * handle;
	T_II_Init			pII_Init;
	T_II_PushRecord		pII_PushRecord;
	T_II_UpdateProgress pII_UpdateProgress;
	T_II_Close			pII_Close;
	T_II_Free			pII_Free;
};

typedef void ( _stdcall * T_PI_Close)(void * handle, bool bHasErrors);
typedef long ( _stdcall * T_PI_PushAllRecords)(void * handle, __int64 nRecordLimit);
typedef long ( _stdcall * T_PI_AddIncomingConnection)(void * handle,
    void * pIncomingConnectionType,
    void * pIncomingConnectionName,
    struct IncomingConnectionInterface *r_IncConnInt);
typedef long ( _stdcall * T_PI_AddOutgoingConnection)(void * handle,
              void * pOutgoingConnectionName,
              struct IncomingConnectionInterface *pIncConnInt);

struct PluginInterface
{
	int								sizeof_PluginInterface;
	void *							handle;
	T_PI_Close						pPI_Close;
	T_PI_AddIncomingConnection		pPI_AddIncomingConnection;
	T_PI_AddOutgoingConnection		pPI_AddOutgoingConnection;
	T_PI_PushAllRecords				pPI_PushAllRecords;
};

struct PresortConnectionInterface;

// Engine definitions

typedef void AlteryxThreadProc(void *pData);
struct PreSortConnectionInterface;
typedef long ( _stdcall * OutputToolProgress)(void * handle, int nToolID, double dPercentProgress);
typedef long ( _stdcall * OutputMessage)(void * handle, int nToolID, int nStatus, wchar_t *pMessage);
typedef unsigned ( _stdcall * BrowseEverywhereReserveAnchor)(void * handle, int nToolId);
typedef struct IncomingConnectionInterface* ( _stdcall * BrowseEverywhereGetII)(void * handle, unsigned nReservationId,  int nToolId, wchar_t * strOutputName);
typedef wchar_t * ( _stdcall * CreateTempFileName)(void * handle, wchar_t * pExt);
typedef long ( _stdcall * PreSort)(void * handle, int nToolId, wchar_t * pSortInfo, struct IncomingConnectionInterface *pOrigIncConnInt, struct IncomingConnectionInterface ** r_ppNewIncConnInt, struct PreSortConnectionInterface ** r_ppPreSortConnInt);
typedef wchar_t * (_stdcall * GetInitVar)(void * handle, wchar_t *pVar);
typedef wchar_t * (_stdcall * GetInitVar2)(void * handle, int nToolId, wchar_t *pVar);

struct EngineInterface {
    int sizeof_EngineInterface;
    void * handle;

    OutputToolProgress pOutputToolProgress;
    OutputMessage pOutputMessage;
    void * pAllocateMemory;
    void * pFreeMemory;
    PreSort pPreSort;
    GetInitVar pGetInitVar;
    CreateTempFileName pCreateTempFileName;
    void * pQueueThread;

    void * pCreateTempFileName2;
    void * pIsLicensed;
    void * pGetConstant;

    GetInitVar2 pGetInitVar2;
    void * pUnlicensedToolCancelled;

    void * pGetConstant2;

    BrowseEverywhereReserveAnchor pBrowseEverywhereReserveAnchor;
    BrowseEverywhereGetII pBrowseEverywhereGetII;

    void * pProfileSetTool;
};

struct PreSortConnectionInterface;


// For the glue


// Plugin methods
struct IncomingConnectionInfo {
    void *handle;
    void *presortString;
    int  cacheSize;
};
void c_configurePlugin(void * handle, struct PluginInterface * pluginInterface);
long c_piPushAllRecords(void * handle, __int64 nRecordLimit);
long go_piPushAllRecords(void * handle, __int64 nRecordLimit);
void c_piClose(void * handle, bool bHasErrors);
void go_piClose(void * handle, bool bHasErrors);
long c_piAddIncomingConnection(void * handle, void * connectionType, void * connectionName, struct IncomingConnectionInterface * incomingInterface);
struct IncomingConnectionInfo *go_piAddIncomingConnection(void * handle, void * connectionType, void * connectionName);
long c_piAddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt);
long go_piAddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt);
struct IncomingConnectionInfo *newSortedIncomingConnectionInfo(void * handle, void * presortString, int cacheSize);
struct IncomingConnectionInfo *newUnsortedIncomingConnectionInfo(void * handle, int cacheSize);

// Incoming interface methods
struct IncomingRecordCache
{
    void*      (*buffer)[];
    int        (*bufferSizes)[];
    int        currentBufferIndex;
    int        recordCount;
    int        recordsInBuffer;
};
struct IncomingConnectionInterface* newIi(void * iiHandle);
void * getIiIndex();
void saveIncomingInterfaceFixedSize(void * handle, int fixedSize);
long c_iiInit(void * handle, void * recordInfoIn);
long go_iiInit(void * handle, void * recordInfoIn);
long c_iiPushRecordCache(void * handle, void * record);
long c_iiPushRecord(void * handle, void * record);
long go_iiPushRecordCache(void * handle, void * cache, int cacheSize);
long go_iiPushRecord(void * handle, void * record);
void c_iiUpdateProgress(void * handle, double percent);
void go_iiUpdateProgress(void * handle, double percent);
void c_iiClose(void * handle);
void go_iiClose(void * handle);
void c_iiFree(void * handle);

// Output connection methods
long c_outputInit(struct IncomingConnectionInterface * connection, void * recordMetaInfoXml);
long c_outputPushRecord(struct IncomingConnectionInterface * connection, void * record);
long c_outputClose(struct IncomingConnectionInterface * connection);
void c_outputUpdateProgress(struct IncomingConnectionInterface * connection, double percent);

// Engine methods
void c_setEngine(struct EngineInterface *pEngineInterface);
void callEngineOutputMessage(int toolId, int status, void * message);
void * callEngineCreateTempFileName(void * ext);
unsigned callEngineBrowseEverywhereReserveAnchor(int toolId);
long callEngineOutputToolProgress(int toolId, double dPercentProgress);
struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(unsigned browseEverywhereAnchorId, int toolId, void * name);
void * callEngineGetInitVar(void * initVar);
void * callEngineGetInitVar2(int toolId, void * initVar);

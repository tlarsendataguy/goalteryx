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
	long ( _stdcall * pII_Init)(void * handle, void * pXmlRecordMetaInfo);
	//T_II_Init			pII_Init;
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
typedef struct IncomingConnectionInterface* ( _stdcall * BrowseEverywhereGetII)(void * handle, unsigned nReservationId,  int nToolId, void * strOutputName);
typedef void * ( _stdcall * CreateTempFileName)(void * handle, void * pExt);
typedef long ( _stdcall * PreSort)(void * handle, int nToolId, void * pSortInfo, struct IncomingConnectionInterface *pOrigIncConnInt, struct IncomingConnectionInterface ** r_ppNewIncConnInt, struct PreSortConnectionInterface ** r_ppPreSortConnInt);

struct EngineInterface {
    int sizeof_EngineInterface;

    void * handle;

    OutputToolProgress pOutputToolProgress;
    OutputMessage pOutputMessage;
    BrowseEverywhereReserveAnchor pBrowseEverywhereReserveAnchor;
    BrowseEverywhereGetII pBrowseEverywhereGetII;
    CreateTempFileName pCreateTempFileName;
    PreSort pPreSort;
};

struct PreSortConnectionInterface;


// For the glue


// Plugin methods
void c_configurePlugin(void * handle, struct PluginInterface * pluginInterface, struct EngineInterface * pluginEngine);
long c_piPushAllRecords(void * handle, __int64 nRecordLimit);
long go_piPushAllRecords(void * handle, __int64 nRecordLimit);
void c_piClose(void * handle, bool bHasErrors);
void go_piClose(void * handle, bool bHasErrors);
long c_piAddIncomingConnection(void * handle, void * connectionType, void * connectionName, struct IncomingConnectionInterface * incomingInterface);
void* go_piAddIncomingConnection(void * handle, void * connectionType, void * connectionName);
long c_piAddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt);
long go_piAddOutgoingConnection(void * handle, void * pOutgoingConnectionName, struct IncomingConnectionInterface *pIncConnInt);

// Incoming interface methods
struct IncomingRecordCache
{
    void*      buffer[10];
    int        bufferSizes[10];
    int        currentBufferIndex;
    int        recordCount;
};
struct IncomingConnectionInterface* newIi(void * iiHandle);
void * getIiIndex();
void saveIncomingInterfaceFixedSize(void * handle, int fixedSize);
long c_iiInit(void * handle, void * recordInfoIn);
long go_iiInit(void * handle, void * recordInfoIn);
long c_iiPushRecord(void * handle, void * record);
long go_iiPushRecordCache(void * handle, void * cache, int cacheSize);
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
void callEngineOutputMessage(struct EngineInterface *pEngineInterface, int toolId, int status, void * message);
void * callEngineCreateTempFileName(struct EngineInterface *pEngineInterface, void * ext);
unsigned callEngineBrowseEverywhereReserveAnchor(struct EngineInterface *pEngineInterface, int toolId);
long callEngineOutputToolProgress(struct EngineInterface *pEngineInterface, int toolId, double dPercentProgress);
struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(struct EngineInterface *pEngineInterface, unsigned browseEverywhereAnchorId, int toolId, void * name);

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

// Engine definitions

typedef void AlteryxThreadProc(void *pData);
struct PreSortConnectionInterface;
typedef long ( _stdcall * OutputToolProgress)(void * handle, int nToolID, double dPercentProgress);
typedef long ( _stdcall * OutputMessage)(void * handle, int nToolID, int nStatus, wchar_t *pMessage);
typedef unsigned ( _stdcall * BrowseEverywhereReserveAnchor)(void * handle, int nToolId);
typedef struct IncomingConnectionInterface* ( _stdcall * BrowseEverywhereGetII)(void * handle, unsigned nReservationId,  int nToolId, void * strOutputName);
typedef void * ( _stdcall * CreateTempFileName)(void * handle, void * pExt);

struct EngineInterface {
    int sizeof_EngineInterface;

    void * handle;

    OutputToolProgress pOutputToolProgress;
    OutputMessage pOutputMessage;
    BrowseEverywhereReserveAnchor pBrowseEverywhereReserveAnchor;
    BrowseEverywhereGetII pBrowseEverywhereGetII;
    CreateTempFileName pCreateTempFileName;
};

// For the glue

void * getPlugin();
typedef long (*outputFunc)(int nToolID, int nStatus, void * pMessage);
void callEngineOutputMessage(struct EngineInterface *pEngineInterface, int toolId, int status, void * message);
unsigned callEngineBrowseEverywhereReserveAnchor(struct EngineInterface *pEngineInterface, int toolId);
struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(struct EngineInterface *pEngineInterface, unsigned browseEverywhereAnchorId, int toolId, void * name);
long callEngineOutputToolProgress(struct EngineInterface *pEngineInterface, int toolId, double dPercentProgress);

long piPushAllRecords(void * handle, __int64 recordLimit);
void piClose(void * handle, bool hasErrors);
long piAddIncomingConnection(void * handle, void * connectionType, void * connectionName, struct IncomingConnectionInterface * incomingInterface);
long piAddOutgoingConnection(void * handle, void * connectionName, struct IncomingConnectionInterface * incomingInterface);
long iiInit(void * handle, void * recordInfoIn);
long iiPushRecord(void * handle, void * record);
void iiUpdateProgress(void * handle, double percent);
void iiClose(void * handle);
void iiFree(void * handle);

long callInitOutput(struct IncomingConnectionInterface * connection, void * recordMetaInfoXml);
long callPushRecord(struct IncomingConnectionInterface * connection, void * record);
long callCloseOutput(struct IncomingConnectionInterface * connection);
void * callEngineCreateTempFileName(struct EngineInterface *pEngineInterface, void * ext);

struct IncomingConnectionInterface* newIi();

struct IncomingRecordCache
{
    void*      buffer[10];
    int        bufferSizes[10];
    int        currentBufferIndex;
    int        recordCount;
};

long pushRecordCache(void * handle, void * cache, int cacheSize);
void closeRecordCache(void * handle);
void * getIiIndex();
void saveIncomingInterfaceFixedSize(void * incomingInterface, int index);
void updateProgress(struct IncomingConnectionInterface * connection, double percent);
void freeRecordCache(void * handle);

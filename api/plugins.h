#include <stdbool.h>
#include <stddef.h>

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

struct EngineInterface {
    int sizeof_EngineInterface;

    void * handle;

    OutputToolProgress pOutputToolProgress;
    OutputMessage pOutputMessage;
    BrowseEverywhereReserveAnchor pBrowseEverywhereReserveAnchor;
    BrowseEverywhereGetII pBrowseEverywhereGetII;
};

// For the glue

void * GetPlugin();
typedef long (*outputFunc)(int nToolID, int nStatus, void * pMessage);
void callEngineOutputMessage(struct EngineInterface *pEngineInterface, int toolId, int status, void * message);
unsigned callEngineBrowseEverywhereReserveAnchor(struct EngineInterface *pEngineInterface, int toolId);
struct IncomingConnectionInterface* callEngineBrowseEverywhereGetII(struct EngineInterface *pEngineInterface, unsigned browseEverywhereAnchorId, int toolId, void * name);

long PiPushAllRecords(void * handle, __int64 recordLimit);
void PiClose(void * handle, bool hasErrors);
long PiAddIncomingConnection(void * handle, void * connectionType, void * connectionName, struct IncomingConnectionInterface * incomingInterface);
long PiAddOutgoingConnection(void * handle, void * connectionName, struct IncomingConnectionInterface * incomingInterface);
long IiInit(void * handle, void * recordInfoIn);
long IiPushRecord(void * handle, void * record);
void IiUpdateProgress(void * handle, double percent);
void IiClose(void * handle);
void IiFree(void * handle);

long callInitOutput(struct IncomingConnectionInterface * connection, void * recordMetaInfoXml);
long callPushRecord(struct IncomingConnectionInterface * connection, void * record);

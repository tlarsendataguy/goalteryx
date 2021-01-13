// Plugin definitions

struct RecordData
{

};

typedef long (* T_II_Init)(void * handle, wchar_t * pXmlRecordMetaInfo);
typedef long (* T_II_PushRecord)(void * handle, char * pRecord);
typedef void (* T_II_UpdateProgress)(void * handle, double dPercent);
typedef void (* T_II_Close)(void * handle);
typedef void (* T_II_Free)(void * handle);

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

typedef void (* T_PI_Close)(void * handle, bool bHasErrors);
typedef long (* T_PI_PushAllRecords)(void * handle, int64_t nRecordLimit);
typedef long (* T_PI_AddIncomingConnection)(void * handle,
    wchar_t * pIncomingConnectionType,
    wchar_t * pIncomingConnectionName,
    struct IncomingConnectionInterface *r_IncConnInt);
typedef long (* T_PI_AddOutgoingConnection)(void * handle,
    wchar_t * pOutgoingConnectionName,
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
typedef long (* OutputToolProgress)(void * handle, int nToolID, double dPercentProgress);
typedef long (* OutputMessage)(void * handle, int nToolID, int nStatus, wchar_t *pMessage);
typedef unsigned (* BrowseEverywhereReserveAnchor)(void * handle, int nToolId);
typedef struct IncomingConnectionInterface* (* BrowseEverywhereGetII)(void * handle, unsigned nReservationId,  int nToolId, wchar_t * strOutputName);
typedef wchar_t * (* CreateTempFileName)(void * handle, wchar_t * pExt);
typedef long (* PreSort)(void * handle, int nToolId, wchar_t * pSortInfo, struct IncomingConnectionInterface *pOrigIncConnInt, struct IncomingConnectionInterface ** r_ppNewIncConnInt, struct PreSortConnectionInterface ** r_ppPreSortConnInt);
typedef wchar_t * (* GetInitVar)(void * handle, wchar_t *pVar);
typedef wchar_t * (* GetInitVar2)(void * handle, int nToolId, wchar_t *pVar);

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

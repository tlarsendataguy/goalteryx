// Plugin definitions

struct RecordData
{

};

typedef long ( _stdcall * T_II_Init)(void * handle, volatile void * pXmlRecordMetaInfo);
typedef long ( _stdcall * T_II_PushRecord)(void * handle, volatile void * pRecord);
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
typedef long ( _stdcall * OutputMessage)(void * handle, int nToolID, int nStatus, volatile wchar_t *pMessage);
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

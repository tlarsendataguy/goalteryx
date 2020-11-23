#include "sdk.h"

const int cacheSize = 4194304; //4mb

/*
** The structure of a plugin handle looks like this:
**
** (struct PluginSharedMemory)
**     toolId (int)
**     toolConfig (void *)
**     engine (struct EngineInterface*)
**     outputAnchors (struct OutputAnchor*)
**         name (void *)
**         isOpen (uint32_t)
**         firstChild (struct OutputConn*)
**             isOpen (uint32_t)
**             ii (struct IncomingInterface*)
**             nextConnection (struct OutputConn*)
**         nextAnchor (struct OutputAnchor*)
**         recordCache (char)
**     totalInputAnchors (uint32_t)
**     closedInputAnchors (uint32_t)
**     inputAnchors (struct InputAnchor*)
**         name (void *)
**         type (void *)
**         isOpen (uint32_t)
**         totalConnections (uint32_t)
**         closedConnections (uint32_t)
**         firstChild (struct InputConnection*)
**             isOpen (uint32_t)
**             metadata (void *)
**             percent (double)
**             nextConnection (struct InputConnection*)
**             recordCache (char[])
**         nextAnchor (struct InputAnchor*)
*/

struct InputConnection {
    uint32_t                isOpen;
    void*                   metadata;
    double                  percent;
    struct InputConnection* nextConnection;
    char                    recordCache[4194304];
};

struct InputAnchor {
    void*                   name;
    void*                   type;
    uint32_t                isOpen;
    uint32_t                totalConnections;
    uint32_t                closedConnections;
    struct InputConnection* firstChild;
    struct InputAnchor*     nextAnchor;
};

struct OutputConn {
    uint32_t                  isOpen;
    struct IncomingInterface* ii;
    struct OutputConn*        nextConnection;
};

struct OutputAnchor {
    void*                name;
    uint32_t             isOpen;
    struct OutputConn*   firstChild;
    struct OutputAnchor* nextAnchor;
    char                 recordCache[4194304];
};

struct PluginSharedMemory {
    uint32_t                toolId;
    void*                   toolConfig;
    struct EngineInterface* engine;
    struct OutputAnchor*    outputAnchors;
    uint32_t                totalInputs;
    uint32_t                closedInputs;
    struct InputAnchor*     inputAnchors;
};

long configurePlugin(int nToolID, void * pXmlProperties, struct EngineInterface *pEngineInterface, struct PluginInterface *r_pluginInterface) {
    return 1;
}
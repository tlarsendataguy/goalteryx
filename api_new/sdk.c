#include "sdk.h"

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
**             record (void *)
**             percent (double)
**             nextConnection (struct InputConnection*)
**         nextAnchor (struct InputAnchor*)
*/

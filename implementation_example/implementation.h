long __declspec(dllexport) PluginEntry(int nToolID,
	void * pXmlProperties,
	void *pEngineInterface,
	void *r_pluginInterface);

long __declspec(dllexport) PluginPresortEntry(int nToolID,
	void * pXmlProperties,
	void *pEngineInterface,
	void *r_pluginInterface);

long __declspec(dllexport) PluginInputEntry(int nToolID,
	void * pXmlProperties,
	void *pEngineInterface,
	void *r_pluginInterface);

long __declspec(dllexport) PluginNoCacheEntry(int nToolID,
	void * pXmlProperties,
	void *pEngineInterface,
	void *r_pluginInterface);

long __declspec(dllexport) NewApiEntry(int nToolID,
	void * pXmlProperties,
	void *pEngineInterface,
	void *r_pluginInterface);

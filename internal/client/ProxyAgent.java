import java.lang.instrument.Instrumentation;

public class ProxyAgent {
    public static void agentmain(String agentArgs, Instrumentation inst) {
        try {
            String[] args = agentArgs.split(":");
            if (args.length != 2) return;
            
            String host = args[0];
            String port = args[1];
            
            System.setProperty("http.proxyHost", host);
            System.setProperty("http.proxyPort", port);
            System.setProperty("https.proxyHost", host);
            System.setProperty("https.proxyPort", port);
            
            // For some JVMs, we also need to disable nonProxyHosts to ensure everything is caught
            System.setProperty("http.nonProxyHosts", "");
            
            System.out.println("[AgentProxy] Successfully injected proxy: " + host + ":" + port);
        } catch (Exception e) {
            System.err.println("[AgentProxy] Error injecting proxy: " + e.getMessage());
        }
    }
}

package com.ctf.challenge.config;

import org.apache.catalina.session.StandardManager;
import org.springframework.boot.web.embedded.tomcat.TomcatServletWebServerFactory;
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TomcatCtfConfig {
    @Bean
    public WebServerFactoryCustomizer<TomcatServletWebServerFactory> strictSessionDeserialization() {
        return factory -> factory.addContextCustomizers(context -> {
            StandardManager manager = new StandardManager();
            manager.setPathname("/tmp/SESSIONS.ser");

            String allowlistRegex = "^(\\[L)?(java\\..*|org\\.springframework\\..*|org\\.apache\\.commons\\.collections\\..*)(;)?$";
            manager.setSessionAttributeValueClassNameFilter(allowlistRegex);

            context.setManager(manager);
        });
    }
}
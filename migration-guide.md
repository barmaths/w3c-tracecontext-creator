\# Migration Guide from v5.x.x into v6.x.x  
  
\### Table of Content:  
1\. <a href="#motivation">Motivation & Rationale Behind v6.x.x.</a>  
2\. <a href="#breaking-changes">Breaking Changes!</a>  
3\. <a href="#new-features">New Features.</a>  
  
\## <a id="motivation"></a>1. Motivation & Rationale Behind v6.x.x.  
\- <a id="otel-limitation">\[1.1\]</a> BEFORE the advent of Open Telemetry Spring Boot Starter, first announced \[here\]([https://opentelemetry.io/blog/2024/spring-starter-stable/),](https://opentelemetry.io/blog/2024/spring-starter-stable/\),) \[Microservice-Starter-Parent Repository\]([https://code.siemens.com/gemini/support/spring-boot-microservice-starter-parent)](https://code.siemens.com/gemini/support/spring-boot-microservice-starter-parent\)) (from now on, `MSPR`, for short) had to use only `@Configuration` classes together with `@EnableNgmXXX` meta annotations for `@ComponentScan` to supply `ngm-default` implementations for all the required/optional `@Bean`s/`@Component`s. This limitation had prohibited the proper use of `@AutoConfiguration`s for company internal custom sub-starters defined under `MSPR`. That, in turn, caused the following items to be accounted for by the downstream services:  
\- <a id="manual-annotations-properties">\[1.1.1\]</a> It was a must to use `@EnableNgmXXX` annotations, besides adding the required sub-starter module into `pom.xml` files, together with additional such `spring.config.import=optional:classpath:<sub-starter-specific-application.yml>` properties in `application.(yml/properties)` files.  
\- <a id="component-scan-order">\[1.1.2\]</a> Spring's component scan order does not guarantee the right ordering of components if extra care is not taken by using explicit `@DependsOn` or other ordering annotations. Otherwise, if your configuration depended upon another bean being created/registered first, you might have been run into issues in the past even if you had `@Conditional`s in place, such as `@Conditional(On/Missing)Bean` annotations.  
\- <a id="component-scan-performance">\[1.1.3\]</a> Even though the explicit component scanning on explicit packages are more visible compared to `@AutoConfiguration-magic`, with every additional “mandatory” sub-starter, the context boot-up time degraded faster with the previous major release.  
  
\- <a id="enforce-mandatory-starters">\[1.2\]</a> BEFORE this new major version, even though classified as mandatory or optional either in some documentation or `ADR`s or some meetings, all the sub-starters defined under `MSPR` were just “technically optional.” If, accidentally or intentionally, the sub-starter is not declared in a `pom.xml`file, then the sub-starter is not enabled at all, independent of the sub-starter classified as optional or mandatory. There was a real need for a way to enforce mandatory sub-starter usage in downstream services.  
\- <a id="bugs">\[1.3\]</a> BEFORE this new major version, there were some major/minor, functional/non-functional bugs/problems in multiple sub-starters to be fixed, besides some additional new features to be implemented.  
  
Besides resolving the above issues, use of proper `@AutoConfiguration` classes in custom starters, except for `business-process-starter`, made the new major version more `spring-bootiful`, maintainable, readable, predictable, performant and pluggable. The following sections elaborates more on the issues, and the breaking changes that need to be applied to resolve them.  
  
\*_NOTE:_\* There are no changes at all other than the version difference of the maven artifact in `business-process-starter:v6.x.x` module compared to its counterpart `business-process-starter:v5.x.x`, regarding the sections of this guide. Watch for future releases of `MSPR` for the similar set of changes for this specific submodule.  
  
\## <a id="breaking-changes"></a>2. Breaking Changes!  
\- <a id="2.1">\[2.1\]</a> To resolve <a href="#enforce-mandatory-starters">\[1.2\]</a> above and classify the sub-starter modules as `imperative` versus `reactive` one should change the `<parent>` reference into either of the following, (replacing the previous `<artifactId>spring-boot-microservice-starter-parent</artifactId>`) depending on the nature of the application:  
\`\`\`_xml  
_<parent>  
<groupId>com.siemens.sbmsp</groupId>  
<artifactId>ngm-imperative-spring-boot-starter</artifactId>  
<version>6.x.x</version>  
</parent>  
\`\`\`  
Or  
\`\`\`_xml  
_<parent>  
<groupId>com.siemens.sbmsp</groupId>  
<artifactId>ngm-reactive-spring-boot-starter</artifactId>  
<version>6.x.x</version>  
</parent>  
\`\`\`  
\- <a id="2.2">\[2.2\]</a> Applying the change in <a href="#2.1">\[2.1\]</a> enforces a different set of mandatory sub-starters to be auto-included, look at the corresponding `pom.xml`files of the maven submodules named, `ngm-imperative-spring-boot-starter` and `ngm-reactive-spring-boot-starter` in `MSPR`, to see these mandatory sub-starters. For instance, referring to `ngm-imperative-spring-boot-starter` results in the auto-inclusion of the following set of maven artifacts:  
\- `logging-starter`  
\- `security-starter`  
\- `health-starter`  
\- `resilience-starter`  
\- `prometheus-metric-starter`  
  
Therefore, if one microservice, before applying the changes mentioned here, doesn't refer to some of these “technically optional but announced as mandatory before” artifacts, it is possible that the autoconfigured beans of the corresponding `@AutoConfiguration` class might result in compile/runtime issues. This could mean the existence of additional beans that did not exist in the `ApplicationContext` before, or this could include clashes of autoconfigured beans with the current `@Configuration/Component` definitions of the microservice itself. To resolve such issues; one can accept the default autoconfigured beans, which are directly defined or referred within `@Import` annotations in the `@AutoConfiguration` classes of the corresponding sub-starter, by removing the existing beans in the microservice itself. Or one can exclude “ALL” the autoconfigured beans by referring to the corresponding `FQN` of `@AutoConfiguration` class in `spring.autoconfigure.exclude` property. Or one can exclude “SOME” of the autoconfigured beans by checking `@ConditionalXXX` annotations over any bean to be excluded and satisfying the negation of that condition. These solutions are applicable for both test and application contexts.  
  
\*_NOTE:_\* Remaining sections of this guide are only applicable for downstream services referring to `ngm-imperative-spring-boot-starter` as their parent. There are no more `sure-to-break`, like <a href="#2.1">\[2.1\]</a>, and `possible-to-break`, like <a href="#2.2">\[2.2\]</a>, changes for downstream services that refer to `ngm-reactive-spring-boot-starter`.  
  
\- <a id="2.3">\[2.3\]</a> To follow the solution to <a href="#manual-annotations-properties">\[1.1.1\]</a> above and prevent compile time failures, one should do the following both:  
\- Remove the `@EnableNgmXXX` annotations from the dedicated `@Configuration` class or from `@SpringBootApplication` class.  
\- Remove `spring.config.import=<optional:>classpath:<sub-starter-specific-application.yml>` properties from `application.(yml/properties)` files.  
  
\- <a id="bug\_refresh">\[2.4\]</a> Before, in `MSPR v5`, secret rotation was handled with `ContextRefresher.refresh()` method, which delegates to `RefreshScope.refreshAll()` method. That, in turn, meant to refresh all the beans annotated with `@RefreshScope` independent of the target beans of the most recent secret rotation i.e., any secret rotation for Db, Redis, Kafka, or IAM client credentials resulted in all the refresh scoped beans to refresh. This bug is resolved in this major release by introducing a new actuator `@WebEndpoint`, with id=`ngmRefresh`, which differentiates target beans of the secret rotation via the sub-path of the input request, rather than using one single `/actuator/refresh` endpoint for all the secret rotations. One needs to reflect this change, by an additional command for Vault to inject, in `deployment.yaml` files of downstream services, like the following mappings show:  
\- `/actuator/ngmRefresh/db` ⇒ PostgreSQL DB secret rotation:  
\`\`\`_yml  
_[vault.hashicorp.com/agent-inject-command-<db\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-command-<db_filename_to_use_in_spring.config.import>:) |  
sh -c '(curl --location --request POST "http://localhost:{{ .Values.port.number }}/actuator/ngmRefresh/db" --header "Content-Type:application/json" --data-raw "{}" -v); exit 0'  
[vault.hashicorp.com/agent-inject-secret-<db\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-secret-<db_filename_to_use_in_spring.config.import>:)  
\# ...  
\`\`\`  
\- `/actuator/ngmRefresh/redis` ⇒ Redis cache secret rotation:  
\`\`\`_yml  
_[vault.hashicorp.com/agent-inject-command-<redis\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-command-<redis_filename_to_use_in_spring.config.import>:) |  
sh -c '(curl --location --request POST "http://localhost:{{ .Values.port.number }}/actuator/ngmRefresh/redis" --header "Content-Type:application/json" --data-raw "{}" -v); exit 0'  
[vault.hashicorp.com/agent-inject-secret-<redis\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-secret-<redis_filename_to_use_in_spring.config.import>:)  
\# ...  
\`\`\`  
\- `/actuator/ngmRefresh/client` ⇒ IAM client secret rotation:  
\`\`\`_yml  
_[vault.hashicorp.com/agent-inject-command-<iam\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-command-<iam_filename_to_use_in_spring.config.import>:) |  
sh -c '(curl --location --request POST "http://localhost:{{ .Values.port.number }}/actuator/ngmRefresh/client" --header "Content-Type:application/json" --data-raw "{}" -v); exit 0'  
[vault.hashicorp.com/agent-inject-secret-<iam\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-secret-<iam_filename_to_use_in_spring.config.import>:)  
\# ...  
\`\`\`  
\- `/actuator/ngmRefresh/kafka` ⇒ Kafka secret rotation:  
\`\`\`_yml  
_[vault.hashicorp.com/agent-inject-command-<kafka\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-command-<kafka_filename_to_use_in_spring.config.import>:) |  
sh -c '(curl --location --request POST "http://localhost:{{ .Values.port.number }}/actuator/ngmRefresh/redis" --header "Content-Type:application/json" --data-raw "{}" -v); exit 0'  
[vault.hashicorp.com/agent-inject-template-<kafka\_filename\_to\_use\_in\_spring.config.import>:](//vault.hashicorp.com/agent-inject-template-<kafka_filename_to_use_in_spring.config.import>:) |  
{{`{{- with secret "internal/data/`}}{{ .Release.Namespace }}/{{ .Values.alm.honeycombName }}/{{ .Values.alm.honeycombName }}-{{ .Release.Name }}{{\`/service-account" -}}  
ngm:  
kafka:  
auth-properties:  
client-secret: {{ .Data.data.secret }}  
{{- end -}}\`}}  
\# ...  
\`\`\`  
\*_NOTES:_\*  
\- `# ...` are used for brevity, one should always use `agent-inject-template-<file_to_be_imported>`, `agent-inject-secret-<file_to_be_imported>` besides `agent-inject-command-<file_to_be_imported>` keys under `spec.template.metadata.annotations` as before in `deployment.yaml` files.  
\- `application.yml` files should also contain the required `spring.config.import` values like so, as before, with the replacement of `<...>` variables in the names of the files injected by Vault:  
\`\`\`_yml  
_spring:  
config:  
import:  
\- optional:file:/vault/secrets/<kafka\_filename\_to\_use\_in\_spring.config.import>  
\- optional:file:/vault/secrets/<db\_filename\_to\_use\_in\_spring.config.import>  
\- optional:file:/vault/secrets/<iam\_filename\_to\_use\_in\_spring.config.import>  
\- optional:file:/vault/secrets/<redis\_filename\_to\_use\_in\_spring.config.import>  
\`\`\`  
\- Note the difference in the value of this key \`[vault.hashicorp.com/agent-inject-template-<kafka\_filename\_to\_use\_in\_spring.config.import>\`,](//vault.hashicorp.com/agent-inject-template-<kafka_filename_to_use_in_spring.config.import>`,) above. This reflects the changes required by the new `ngm.kafka` prefixed `@ConfigurationProperties` of type `NgmKafkaProperties`, which will be revisited in more detail below in the section <a href="#kafka-auth-starter-changes">\[2.8\]</a>.  
  
\- <a id="security-starter-changes">\[2.5\]</a> Changes in `security-starter`:  
  
\- <a id="health-starter-changes">\[2.6\]</a> Changes in `health-starter`:  
  
\- <a id="logging-starter-changes">\[2.7\]</a> Changes in `logging-starter`:  
  
\- <a id="kafka-auth-starter-changes">\[2.8\]</a> Changes in `kafka-auth-starter`:  
  
\## <a id="new-features"></a>3. New Features.  
\- Mention a new module `starter-commons` and its content, including `NgmRefreshableXXX` beans  
\- Mention `NgmRestClientBuilder` bean
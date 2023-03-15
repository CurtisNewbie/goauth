package com.curtisnewbie.goauth.client;

import com.curtisnewbie.common.vo.*;
import lombok.Data;
import lombok.extern.slf4j.*;
import org.springframework.beans.factory.*;
import org.springframework.beans.factory.annotation.*;
import org.springframework.core.env.*;
import org.springframework.util.*;

import java.util.*;
import java.util.concurrent.*;
import java.util.stream.Collectors;

/**
 * Reporter of REST paths parsed by RestPathScanner
 *
 * @author yongj.zhuang
 */
@Slf4j
public class RestPathReporter implements InitializingBean {

    public static final String DISABLE_REPORT_KEY = "goauth.path.report.disabled";

    @Autowired
    private RestPathScanner restPathScanner;
    @Autowired
    private GoAuthClient goAuthClient;
    @Autowired
    private Environment env;

    @Override
    public void afterPropertiesSet() throws Exception {
        final String group = env.getProperty("spring.application.name");

        final boolean disabled = Boolean.parseBoolean(env.getProperty(DISABLE_REPORT_KEY, "false"));
        if (!disabled) {
            restPathScanner.onParsed(restPaths -> {
                final StopWatch sw = new StopWatch();
                sw.start();
                reportResources(restPaths, goAuthClient);
                reportPaths(restPaths, group, goAuthClient);
                sw.stop();
                log.info("GoAuth RestPath Reported, took: {}ms", sw.getTotalTimeMillis());
            });
        }
    }

    protected static void reportResources(List<RestPathScanner.RestPath> restPaths, GoAuthClient goAuthClient) {
        final Map<String /* code */, PResource> resources = restPaths.stream()
                .filter(p -> !p.requestPath.startsWith("/remote"))
                .filter(p -> StringUtils.hasText(p.pathDoc.resCode()))
                .map(p -> new PResource(p.pathDoc.resCode(), p.pathDoc.resName()))
                .collect(Collectors.toMap(r -> r.code, r -> r, (a, b) -> a));

        try {
            resources.forEach((k, v) -> goAuthClient.addResource(new AddResourceReq(v.name, v.code)).assertIsOk());
        } catch (Throwable e) {
            log.error("Failed to report resources to goauth, resources: {}", resources.values(), e);
        }
    }

    protected static void reportPaths(List<RestPathScanner.RestPath> restPaths, String group, GoAuthClient goAuthClient) {
        final List<AddPathReq> reqs = restPaths.stream()
                .filter(p -> !p.requestPath.startsWith("/remote"))
                .map(p -> {
                    final AddPathReq ar = new AddPathReq();
                    ar.setUrl("/" + group + p.getCompletePath());
                    ar.setGroup(group);
                    ar.setType(p.pathDoc.type());
                    ar.setDesc(p.pathDoc.description());
                    ar.setResCode(p.pathDoc.resCode());
                    return ar;
                })
                .collect(Collectors.toList());

        try {
            batchReportPaths(reqs, goAuthClient);
            goAuthClient.reloadPathCache().assertIsOk();
        } catch (Throwable e) {
            log.error("Failed to report path to goauth, reqs: {}", reqs, e);
        }
    }

    protected static void batchReportPaths(List<AddPathReq> reqList, GoAuthClient goAuthClient) {
        final BatchAddPathReq req = new BatchAddPathReq();
        req.setReqs(reqList);
        final Result<Void> res = goAuthClient.batchAddPath(req);
        if (!res.isOk()) {
            log.error("Failed to report path to goauth, reqs: {}, error code: {}, error msg: {}",
                    reqList, res.getErrorCode(), res.getMsg());
            return;
        }
        log.info("Reported {} paths to goauth", reqList.size());
    }

    protected static void reportPath(String group, String url, PathType type, GoAuthClient goAuthClient) {
        try {
            AddPathReq req = new AddPathReq();
            req.setGroup(group);
            req.setType(type);
            req.setUrl(url);

            final Result<Void> res = goAuthClient.addPath(req);
            if (!res.isOk()) {
                log.error("Failed to report path to goauth, group: {}, type: {}, url: {}, error code: {}, error msg: {}",
                        req.getGroup(), req.getType(), req.getUrl(), res.getErrorCode(), res.getMsg());
                return;
            }

            log.info("Reported path '{}' to goauth", req.getUrl());
        } catch (Throwable e) {
            log.error("Failed to report path to goauth, group: {}, type: {}, url: {}", group, type, url, e);
        }
    }

    @Data
    private static class PResource {
        private String code;
        private String name;

        public PResource(String code, String name) {
            this.code = code;
            this.name = name;
        }
    }
}


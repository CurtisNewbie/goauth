package com.curtisnewbie.goauth.client;

import com.curtisnewbie.common.vo.*;
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
                reportPaths(restPaths, group, goAuthClient);
                sw.stop();
                log.info("GoAuth RestPath Reported, took: {}ms", sw.getTotalTimeMillis());
            });
        }
    }

    protected static void reportPaths(List<RestPathScanner.RestPath> restPaths, String group, GoAuthClient goAuthClient) {
        final List<String> paths = restPaths.stream()
                .map(p -> "/" + group + p.getCompletePath())
                .distinct()
                .collect(Collectors.toList());
        batchReportPaths(group, paths, PathType.PROTECTED, goAuthClient);
        goAuthClient.reloadPathCache();
    }

    protected static void batchReportPaths(String group, List<String> urls, PathType type, GoAuthClient goAuthClient) {
        try {
            BatchAddPathReq req = new BatchAddPathReq();
            req.setGroup(group);
            req.setType(type);
            req.setUrls(urls);

            final Result<Void> res = goAuthClient.batchAddPath(req);
            if (!res.isOk()) {
                log.error("Failed to report path to goauth, group: {}, type: {}, url: {}, error code: {}, error msg: {}",
                        req.getGroup(), req.getType(), req.getUrls(), res.getErrorCode(), res.getMsg());
                return;
            }

            log.info("Reported paths '{}' to goauth", req.getUrls());
        } catch (Throwable e) {
            log.error("Failed to report path to goauth, group: {}, type: {}, url: {}", group, type, urls, e);
        }
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
}

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
        final List<AddPathReq> reqs = restPaths.stream()
                .filter(p -> !p.requestPath.startsWith("/remote"))
                .map(p -> {
                    final AddPathReq ar = new AddPathReq();
                    ar.setUrl("/" + group + p.getCompletePath());
                    ar.setGroup(group);
                    ar.setType(PathType.PROTECTED);
                    ar.setDesc(p.description);
                    return ar;
                })
                .collect(Collectors.toList());

        try {
            batchReportPaths(reqs, goAuthClient);
            goAuthClient.reloadPathCache();
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
}

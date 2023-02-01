package com.curtisnewbie.goauth.client;

import com.curtisnewbie.common.vo.*;
import lombok.extern.slf4j.*;
import org.springframework.beans.factory.*;
import org.springframework.beans.factory.annotation.*;
import org.springframework.core.env.*;

/**
 * Reporter of REST paths parsed by RestPathScanner
 *
 * @author yongj.zhuang
 */
@Slf4j
public class RestPathReporter implements InitializingBean {

    @Autowired
    private RestPathScanner restPathScanner;
    @Autowired
    private GoAuthClient goAuthClient;
    @Autowired
    private Environment env;

    @Override
    public void afterPropertiesSet() throws Exception {
        final String group = env.getProperty("spring.application.name");

        restPathScanner.onParsed(restPaths -> {
            restPaths.stream()
                    .map(p -> "/" + group + p.getCompletePath())
                    .distinct()
                    .forEach(url -> {
                        AddPathReq req = new AddPathReq();
                        req.setGroup(group);
                        req.setType(PathType.PROTECTED);
                        req.setUrl(url);
                        final Result<Void> res = goAuthClient.addPath(req);
                        if (!res.isOk()) {
                            log.error("Failed to report path to goauth, group: {}, type: {}, url: {}, error code: {}, error msg: {}",
                                    req.getGroup(), req.getType(), req.getUrl(), res.getErrorCode(), res.getMsg());
                        }
                    });
        });
    }
}

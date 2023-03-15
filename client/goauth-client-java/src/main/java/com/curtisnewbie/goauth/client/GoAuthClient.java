package com.curtisnewbie.goauth.client;

import com.curtisnewbie.common.vo.Result;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

/**
 * @author yongj.zhuang
 */
@FeignClient(value = "goauth", path = "/remote")
public interface GoAuthClient {

    @PostMapping("/path/resource/access-test")
    Result<TestResAccessResp> testResAccess(@RequestBody TestResAccessReq req);

    @PostMapping("/path/add")
    Result<Void> addPath(@RequestBody AddPathReq req);

    @PostMapping("/path/batch/add")
    Result<Void> batchAddPath(@RequestBody BatchAddPathReq req);

    @PostMapping("/role/info")
    Result<RoleInfoResp> getRoleInfo(@RequestBody RoleInfoReq req);

    @PostMapping("/resource/add")
    Result<Void> addResource(@RequestBody AddResourceReq req);

    @PostMapping("/path/cache/reload")
    Result<Void> reloadPathCache();

}

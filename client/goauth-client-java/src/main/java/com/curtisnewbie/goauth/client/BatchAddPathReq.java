package com.curtisnewbie.goauth.client;

import lombok.Data;

import java.util.List;

/**
 * @author yongj.zhuang
 */
@Data
public class BatchAddPathReq {
    private List<AddPathReq> reqs;
}

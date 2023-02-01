package com.curtisnewbie.goauth.client;

import lombok.*;
import lombok.extern.slf4j.*;
import org.junit.jupiter.api.*;
import org.springframework.boot.autoconfigure.*;
import org.springframework.boot.test.context.*;
import org.springframework.context.annotation.*;
import org.springframework.stereotype.*;
import org.springframework.web.bind.annotation.*;

import java.util.*;

/**
 * @author yongj.zhuang
 */
@Slf4j
//@SpringBootTest(classes = RestPathScannerTest.class)
//@SpringBootApplication
public class RestPathScannerTest {

    @Test
    public void should_parse_rest_path() {
        List<RestPathScanner.RestPath> l = new ArrayList<>();
        RestPathScanner.parseRestPath(l, DummyCtrl.class, t -> t);
        System.out.println(l);
        Assertions.assertEquals(12, l.size()); // 8 (for /any) + 4
    }

    @Test
    public void should_get_complete_path() {
        Assertions.assertEquals("/dummy", new RestPathScanner.RestPath("dummy", "", RequestMethod.GET).getCompletePath());
        Assertions.assertEquals("/dummy/info", new RestPathScanner.RestPath("dummy", "info", RequestMethod.GET).getCompletePath());
        Assertions.assertEquals("/info", new RestPathScanner.RestPath("", "/info", RequestMethod.GET).getCompletePath());
        Assertions.assertEquals("/dummy/info", new RestPathScanner.RestPath("dummy", "info/////", RequestMethod.GET).getCompletePath());
    }

    @Component
    @RestController
    @RequestMapping("/dummy")
    public static class DummyCtrl {

        @RequestMapping("/any")
        public void any() {
        }

        @GetMapping("/get-info")
        public void getInfo() {
        }

        @PutMapping("/put-info")
        public void putInfo() {
        }

        @PostMapping("/post-info")
        public void postInfo() {
        }

        @DeleteMapping("/del-info")
        public void deleteInfo() {
        }
    }
}

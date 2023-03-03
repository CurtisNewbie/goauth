package com.curtisnewbie.goauth.client;

import java.lang.annotation.*;

/**
 * Path Documentation
 *
 * @author yongj.zhuang
 */
@Documented
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface PathDoc {

    String description() default "";
}

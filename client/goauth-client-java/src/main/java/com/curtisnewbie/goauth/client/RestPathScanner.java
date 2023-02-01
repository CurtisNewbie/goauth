package com.curtisnewbie.goauth.client;

import lombok.*;
import org.springframework.beans.*;
import org.springframework.context.*;
import org.springframework.stereotype.*;
import org.springframework.web.bind.annotation.*;

import java.lang.annotation.*;
import java.lang.reflect.*;
import java.util.*;
import java.util.function.*;

import static java.util.Collections.*;

/**
 * Scanner of REST Path
 * <p>
 * Potential candidates include beans that are annotated with @Controller and @RestController
 *
 * @author yongj.zhuang
 */
public class RestPathScanner implements ApplicationContextAware {

    private static final Map<Class<? extends Annotation>, MappingPathParser> clz2Parser = new HashMap<>();
    private volatile List<RestPath> parsedRestPaths = new ArrayList<>();
    private List<Consumer<List<RestPath>>> onParsed = new ArrayList<>();

    static {
        clz2Parser.put(RequestMapping.class, o -> {
            RequestMapping rm = (RequestMapping) o;
            RequestMethod method = rm.method().length > 0 ? rm.method()[0] : null;

            String path;
            if (rm.value().length > 0) path = rm.value()[0];
            else path = rm.path().length > 0 ? rm.path()[0] : "";

            if (method != null) {
                return singletonList(new ParsedMapping(path, method));
            }

            // all HTTP methods
            final List<ParsedMapping> parsed = new ArrayList<>();
            for (RequestMethod mtd : RequestMethod.values()) {
                parsed.add(new ParsedMapping(path, mtd));
            }
            return parsed;
        });
        clz2Parser.put(GetMapping.class, o -> {
            GetMapping gm = (GetMapping) o;
            RequestMethod method = RequestMethod.GET;
            if (gm.value().length > 0) return singletonList(new ParsedMapping(gm.value()[0], method));
            return singletonList(new ParsedMapping(gm.path().length > 0 ? gm.path()[0] : "", method));
        });
        clz2Parser.put(PutMapping.class, o -> {
            PutMapping pm = (PutMapping) o;
            RequestMethod method = RequestMethod.PUT;
            if (pm.value().length > 0) return singletonList(new ParsedMapping(pm.value()[0], method));
            return singletonList(new ParsedMapping(pm.path().length > 0 ? pm.path()[0] : "", method));
        });
        clz2Parser.put(PostMapping.class, o -> {
            PostMapping pm = (PostMapping) o;
            RequestMethod method = RequestMethod.POST;
            if (pm.value().length > 0) return singletonList(new ParsedMapping(pm.value()[0], method));
            return singletonList(new ParsedMapping(pm.path().length > 0 ? pm.path()[0] : "", method));
        });
        clz2Parser.put(DeleteMapping.class, o -> {
            DeleteMapping dm = (DeleteMapping) o;
            RequestMethod method = RequestMethod.DELETE;
            if (dm.value().length > 0) return singletonList(new ParsedMapping(dm.value()[0], method));
            return singletonList(new ParsedMapping(dm.path().length > 0 ? dm.path()[0] : "", method));
        });
    }

    @Override
    public void setApplicationContext(ApplicationContext appCtx) throws BeansException {
        final Map<String, Object> beans = appCtx.getBeansWithAnnotation(Controller.class);
        beans.forEach((k, v) -> {
            List<RestPath> restPaths = new ArrayList<>();
            Class<?> beanClz = v.getClass();
            parseRestPath(restPaths, beanClz);

            synchronized (this) {
                this.parsedRestPaths = restPaths;
                if (!this.onParsed.isEmpty()) {
                    this.onParsed.forEach(callback -> {
                        callback.accept(new ArrayList<>(this.parsedRestPaths));
                    });
                }
            }
        });
    }

    /** Register onParsed callback */
    public void onParsed(Consumer<List<RestPath>> callback) {
        if (callback == null) return;
        synchronized (this) {
            if (this.parsedRestPaths != null) {
                callback.accept(new ArrayList<>(this.parsedRestPaths));
            } else {
                this.onParsed.add(callback);
            }
        }
    }

    public static void parseRestPath(List<RestPath> restPathList, Class<?> beanClz) {
        String rootPath = "";
        final RequestMapping rootMapping = beanClz.getDeclaredAnnotation(RequestMapping.class);
        if (rootMapping != null) {
            final List<ParsedMapping> parsed = clz2Parser.get(RequestMapping.class).parsed(rootMapping);
            if (!parsed.isEmpty())
                rootPath = parsed.get(0).requestPath;
        }

        final Method[] methods = beanClz.getDeclaredMethods();
        for (int i = 0; i < methods.length; i++) {
            Method m = methods[i];

            for (Annotation mda : m.getDeclaredAnnotations()) {
                Class<?> typ = mda.annotationType();
                if (clz2Parser.containsKey(typ)) {
                    final List<ParsedMapping> parsed = clz2Parser.get(typ).parsed(mda);
                    for (ParsedMapping pm : parsed) {
                        restPathList.add(new RestPath(rootPath, pm.requestPath, pm.httpMethod));
                    }

                    break; // normally, a method can only have one mapping
                }
            }
        }
    }

    /**
     * Parsed REST Path, thread-safe
     */
    @ToString
    @AllArgsConstructor
    public static class RestPath {
        public final String rootPath;
        public final String requestPath;
        public final RequestMethod httpMethod;

        public String getCompletePath() {
            String rtp = rootPath != null ? rootPath.trim() : "";
            String rqp = requestPath != null ? requestPath.trim() : "";
            if (!rtp.isEmpty() && !rtp.startsWith("/")) rtp = "/" + rtp;
            if (!rqp.isEmpty() && !rqp.startsWith("/")) rqp = "/" + rqp;

            int j = -1;
            for (int i = rqp.length() - 1; i > -1; i--) {
                if (rqp.charAt(i) == '/') j = i;
                else break;
            }
            if (j > -1) rqp = rqp.substring(0, j); // remove trailing '/'


            // not sure about returning "/" like this when both rtp and rqp are empty, doesn't seem like that it will actually happen :D
            /*
                String pt = rtp + rqp;
                return pt.isEmpty() ? "/" : pt;
             */

            return rtp + rqp;
        }
    }

    @FunctionalInterface
    private interface MappingPathParser {
        List<ParsedMapping> parsed(Annotation o);
    }

    @AllArgsConstructor
    public static class ParsedMapping {
        public final String requestPath;
        public final RequestMethod httpMethod;
    }
}

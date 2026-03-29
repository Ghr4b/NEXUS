package com.ctf.challenge.util;

public class SqlValidator {

    private static final String[] HARD_BLOCK_KEYWORDS = {
            "UNION", "SELECT", "INSERT", "UPDATE", "DELETE",
            "ALTER", "DROP", "CREATE", "EXEC", "EXECUTE",
            "MERGE", "REPLACE", "TRUNCATE", "GRANT", "REVOKE",
            "CALL", "INFORMATION_SCHEMA", "PG_SLEEP", "SLEEP",
            "BENCHMARK", "LOAD_FILE", "OUTFILE", "DUMPFILE",
            "WAITFOR", "DBMS_",

            // H2 RCE primitives
            "RUNSCRIPT", "LINK_SCHEMA", "LINK_TABLE",
            "TRACE_LEVEL", "INIT=",

            // H2 aliases / alternate names for SCRIPT
            "SCRIPT",

            // Encoding / type conversion (unintended payload delivery)
            "HEXTORAW", "RAWTOHEX",
            "TO_CHAR", "TO_NUMBER", "TO_DATE",
            "CONVERT", "CAST",
            "DECODE", "ENCODE",
            "ASCII", "ORD", "BIN", "OCT",
            "CHAR", "NCHAR", "CHR", // all aliases for character-from-int
            "FROMCHARCODE",

            // Heap OOM alternatives to REPEAT
            "LPAD", "RPAD", "SPACE",

            // String manipulation that can reconstruct blocked keywords
            "CONCAT", // alternative to ||
            "INSERT", // also a string function in H2
            "OVERLAY", // SQL string replacement
            "SUBSTRING", "SUBSTR", "MID", // can slice/reconstruct strings
            "REVERSE",
            "REPLACE", // can reconstruct blocked tokens
            "TRANSLATE",

            // Math/system functions with side effects
            "RAND", "RANDOM", // timing side-channel
            "HASH",
            "ENCRYPT", "DECRYPT",
            "COMPRESS", "DECOMPRESS",

            // Set operators
            "EXCEPT", "MINUS", "INTERSECT",

            // XSS hard block
            "<SCRIPT", "JAVASCRIPT:", "ONERROR=", "ONLOAD=",
            "ONCLICK=", "ALERT", "CONFIRM", "PROMPT",
            "DOCUMENT.COOKIE", "DOCUMENT.WRITE", "EVAL",
            "VBSCRIPT:", "EXPRESSION",

            // Comment sequences
            "--", "/*", "*/", "#",

            // Statement terminator
            ";",
    };

    public static String sanitize(String input) {
        if (input == null || input.trim().isEmpty()) {
            return "";
        }

        String upper = input.toUpperCase();
        for (String kw : HARD_BLOCK_KEYWORDS) {
            if (upper.contains(kw)) {
                throw new IllegalArgumentException(
                        "Your search contained potentially harmful characters and was blocked. " +
                                "Please use plain product names or descriptions.");
            }
        }

        String sanitized = input;
        sanitized = sanitized.replaceAll("(?i)FILE_WRITE", "");
        sanitized = sanitized.replaceAll("(?i)FILE_READ", "");
        sanitized = sanitized.replaceAll("(?i)CSVWRITE", "");
        sanitized = sanitized.replaceAll("(?i)CSVREAD", "");
        sanitized = sanitized.replaceAll("(?i)\\bEXCEPT\\b", "");
        sanitized = sanitized.replaceAll("(?i)\\bMINUS\\b", "");
        sanitized = sanitized.replaceAll("(?i)\\bINTERSECT\\b", "");
        sanitized = sanitized.replaceAll("(?i)\\bWAITFOR\\b", "");
        sanitized = sanitized.replaceAll("(?i)<[a-zA-Z/!][^>]*>", "");
        sanitized = sanitized.replaceAll("(?i)&[a-zA-Z]{2,6};", "");

        return sanitized;
    }
}

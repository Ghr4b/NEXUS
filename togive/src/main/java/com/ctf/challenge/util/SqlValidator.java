package com.ctf.challenge.util;

public class SqlValidator {

    private static final String[] HARD_BLOCK_KEYWORDS = {
            "UNION", "SELECT", "INSERT", "UPDATE", "DELETE",
            "ALTER", "DROP", "CREATE", "EXEC", "EXECUTE",
            "MERGE", "REPLACE", "TRUNCATE", "GRANT", "REVOKE",
            "CALL", "INFORMATION_SCHEMA", "PG_SLEEP", "SLEEP",
            "BENCHMARK", "LOAD_FILE", "OUTFILE", "DUMPFILE",
            "WAITFOR", "DBMS_",

            "RUNSCRIPT", "LINK_SCHEMA", "LINK_TABLE",
            "TRACE_LEVEL", "INIT=",

            "SCRIPT",

            "HEXTORAW", "RAWTOHEX",
            "TO_CHAR", "TO_NUMBER", "TO_DATE",
            "CONVERT", "CAST",
            "DECODE", "ENCODE",
            "ASCII", "ORD", "BIN", "OCT",
            "CHAR", "NCHAR", "CHR",
            "FROMCHARCODE",

            "LPAD", "RPAD", "SPACE",


            "CONCAT",
            "INSERT",
            "OVERLAY",
            "SUBSTRING", "SUBSTR", "MID",
            "REVERSE",
            "REPLACE",
            "TRANSLATE","FILE_READ","CSVREAD","CSVWRITE"

            "RAND", "RANDOM",
            "HASH",
            "ENCRYPT", "DECRYPT",
            "COMPRESS", "DECOMPRESS",

            "EXCEPT", "MINUS", "INTERSECT",

            "<SCRIPT", "JAVASCRIPT:", "ONERROR=", "ONLOAD=",
            "ONCLICK=", "ALERT", "CONFIRM", "PROMPT",
            "DOCUMENT.COOKIE", "DOCUMENT.WRITE", "EVAL",
            "VBSCRIPT:", "EXPRESSION",

            "--", "/*", "*/", "#",

            ";",
    };

    public static String sanitize(String input) {
        if (input == null || input.trim().isEmpty()) {
            return "";
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

        String upper = sanitized.toUpperCase();
        for (String kw : HARD_BLOCK_KEYWORDS) {
            if (upper.contains(kw)) {
                throw new IllegalArgumentException(
                        "Your search contained potentially harmful characters and was blocked. " +
                                "Please use plain product names or descriptions.");
            }
        }
        return sanitized;
    }
}

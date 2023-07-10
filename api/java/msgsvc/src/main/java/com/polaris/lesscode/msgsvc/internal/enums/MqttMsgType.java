package com.polaris.lesscode.msgsvc.internal.enums;

public enum MqttMsgType {

    PROJECT("PRO", "项目"),
    ISSUE("ISSUE", "任务"),
    TAG("TAG", "标签"),
    MEMBER("MEMBER", "成员"),
    FORM_CONFIG("FORM_CONFIG", "表单配置"),
    VIEW("VIEW", "视图"),

    ;

    private final String code;

    private final String desc;

    MqttMsgType(String code, String desc) {
        this.code = code;
        this.desc = desc;
    }

    public String getCode() {
        return code;
    }

    public String getDesc() {
        return desc;
    }
}

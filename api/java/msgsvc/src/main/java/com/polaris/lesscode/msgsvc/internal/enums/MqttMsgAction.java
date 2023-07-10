package com.polaris.lesscode.msgsvc.internal.enums;

public enum MqttMsgAction {

    ADD("ADD", "添加"),
    DEL("DEL", "删除"),
    MODIFY("MODIFY", "修改"),
    MODIFYSORT("MODIFYSORT", "修改排序"),
    MOVE("MOVE", "移动"),
    Archive("Archive", "归档"),
    CancelArchive("CancelArchive", "取消归档"),

    ;

    private final String code;

    private final String desc;

    MqttMsgAction(String code, String desc) {
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

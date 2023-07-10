package com.polaris.lesscode.project.internal.req;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class AppEvent {
    @JsonProperty("orgId")
    private Long orgId;

    @JsonFormat(shape = JsonFormat.Shape.STRING)
    @JsonProperty("appId")
    private Long appId;

    @JsonProperty("projectId")
    private Long projectId;

    @JsonProperty("userId")
    private Long userId;

    @JsonProperty("app")
    private Object app;

    @JsonProperty("project")
    private Object project;

    @JsonProperty("chat")
    private Object chat;
}

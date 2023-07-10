package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

@Data
public class ProjectDetail {

    private String notice;

    private Integer isEnableWorkHours;

    private Integer isSyncOutCalendar;

    private String fsChatSettings;
}

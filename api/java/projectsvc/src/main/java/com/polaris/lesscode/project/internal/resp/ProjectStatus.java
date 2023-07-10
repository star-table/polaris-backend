package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

@Data
public class ProjectStatus {

    private Long id;

    private String name;

    private Integer type;

    private String fontStyle;

    private String bgStyle;

    private Integer sort;

	private Integer isInitStatus;

	private Integer category;
}

package com.polaris.lesscode.msgsvc.internal.req;

import lombok.Data;

import java.util.List;

@Data
public class PartialRefresh {

    private long ObjectId;

    private List<ObjectField> fields;

}

package com.polaris.lesscode.msgsvc.internal.req;

import com.sun.corba.se.spi.ior.ObjectId;
import lombok.Data;

@Data
public class ObjectField {

    private String field;

    private Object value;

    private Object newObjectValue;
}

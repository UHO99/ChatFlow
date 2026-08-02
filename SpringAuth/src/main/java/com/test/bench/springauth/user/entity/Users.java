package com.test.bench.springauth.user.entity;

import com.test.bench.springauth.user.dto.UserDto;
import jakarta.persistence.*;
import lombok.Builder;
import lombok.Getter;

import java.time.OffsetDateTime;

@Getter
@Entity
@Builder
public class Users {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private long id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    private String username;

    @Column(nullable = false, unique = true)
    private String email;

    @Column(nullable = false)
    private String password;

    @Column(nullable = false, updatable = false, insertable = false)
    private OffsetDateTime createdAt;

    public UserDto toDto(){
        return UserDto.builder()
                .name(this.name)
                .username(this.username)
                .email(this.email)
                .build();
    }

}

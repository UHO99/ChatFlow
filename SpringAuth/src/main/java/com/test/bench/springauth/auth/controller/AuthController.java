package com.test.bench.springauth.auth.controller;

import com.test.bench.springauth.auth.dto.LoginRequestDto;
import com.test.bench.springauth.auth.dto.TokenResponseDto;
import com.test.bench.springauth.auth.service.AuthService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
public class AuthController {

    private static final String BEARER_PREFIX = "Bearer ";

    private final AuthService authService;

    @PostMapping("/login")
    public ResponseEntity<TokenResponseDto> login(@RequestBody LoginRequestDto loginRequestDto) {
        return ResponseEntity.ok(authService.login(loginRequestDto));
    }

    @PostMapping("/logout")
    public ResponseEntity<Void> logout(@RequestHeader(HttpHeaders.AUTHORIZATION) String authorization) {
        authService.logout(resolveToken(authorization));
        return ResponseEntity.ok().build();
    }

    private String resolveToken(String authorization) {
        return authorization.startsWith(BEARER_PREFIX) ? authorization.substring(BEARER_PREFIX.length()) : authorization;
    }

}

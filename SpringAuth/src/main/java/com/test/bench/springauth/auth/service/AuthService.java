package com.test.bench.springauth.auth.service;

import com.test.bench.springauth.auth.dto.LoginRequestDto;
import com.test.bench.springauth.auth.dto.TokenResponseDto;

public interface AuthService {

    TokenResponseDto login(LoginRequestDto loginRequestDto);

    void logout(String token);

}

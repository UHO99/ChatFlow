package com.test.bench.springauth.auth.service;

import com.test.bench.springauth.auth.dto.LoginRequestDto;
import com.test.bench.springauth.auth.dto.TokenResponseDto;
import com.test.bench.springauth.auth.jwt.JwtTokenProvider;
import com.test.bench.springauth.auth.jwt.TokenBlacklistService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

    private static final String TOKEN_TYPE = "Bearer";

    private final AuthenticationManager authenticationManager;
    private final JwtTokenProvider jwtTokenProvider;
    private final TokenBlacklistService tokenBlacklistService;

    @Override
    public TokenResponseDto login(LoginRequestDto loginRequestDto) {
        Authentication authentication = authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(loginRequestDto.getUsername(), loginRequestDto.getPassword()));

        String token = jwtTokenProvider.generateToken(authentication.getName());
        long expiresIn = jwtTokenProvider.getExpiration(token);

        return TokenResponseDto.builder()
                .accessToken(token)
                .tokenType(TOKEN_TYPE)
                .expiresIn(expiresIn)
                .build();
    }

    @Override
    public void logout(String token) {
        long ttl = jwtTokenProvider.getExpiration(token);
        if (ttl > 0) {
            tokenBlacklistService.blacklist(token, ttl);
        }
    }

}

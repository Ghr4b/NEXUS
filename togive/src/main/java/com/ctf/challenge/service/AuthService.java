package com.ctf.challenge.service;

import com.ctf.challenge.model.PromoCode;
import com.ctf.challenge.model.User;
import com.ctf.challenge.repository.PromoRepository;
import com.ctf.challenge.repository.UserRepository;
import jakarta.servlet.http.HttpSession;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class AuthService {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private PromoRepository promoRepository;

    public User login(String username, String password) {
        User user = userRepository.findByUsername(username);
        if (user != null && user.getPassword().equals(password)) {
            return user;
        }
        return null;
    }

    public User register(String username, String password) {
        if (userRepository.findByUsername(username) != null) {
            return null;
        }
        User user = new User();
        user.setUsername(username);
        user.setPassword(password);
        user.setRole("USER");
        userRepository.save(user);

        User savedUser = userRepository.findByUsername(username);

        PromoCode promo = new PromoCode();
        promo.setCode("CYBERSHOP-" + UUID.randomUUID().toString().substring(0, 8).toUpperCase());
        promo.setDiscountPercentage(15);
        promo.setIsActive(true);
        promo.setUserId(savedUser.getId());
        promoRepository.save(promo);

        return savedUser;
    }

    public void setSession(HttpSession session, User user) {
        session.setAttribute("USER_ID", user.getId());
        session.setAttribute("USERNAME", user.getUsername());
        session.setAttribute("ROLE", user.getRole());
    }

    public void clearSession(HttpSession session) {
        session.invalidate();
    }
}

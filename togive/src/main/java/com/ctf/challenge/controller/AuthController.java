package com.ctf.challenge.controller;

import com.ctf.challenge.model.User;
import com.ctf.challenge.service.AuthService;
import jakarta.servlet.http.HttpSession;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@Controller
public class AuthController {

    @Autowired
    private AuthService authService;

    @GetMapping("/login")
    public String showLoginFrame() {
        return "auth/login";
    }

    @PostMapping("/login")
    public String performLogin(@RequestParam String username, @RequestParam String password, HttpSession session,
            Model model) {
        User user = authService.login(username, password);
        if (user != null) {
            authService.setSession(session, user);
            return "redirect:/";
        }
        model.addAttribute("error", "Invalid credentials");
        return "auth/login";
    }

    @GetMapping("/register")
    public String showRegisterFrame() {
        return "auth/register";
    }

    @PostMapping("/register")
    public String performRegister(@RequestParam String username, @RequestParam String password, HttpSession session,
            Model model) {
        User user = authService.register(username, password);
        if (user != null) {
            authService.setSession(session, user);
            return "redirect:/";
        }
        model.addAttribute("error", "Username already taken");
        return "auth/register";
    }

    @GetMapping("/logout")
    public String logout(HttpSession session) {
        authService.clearSession(session);
        return "redirect:/";
    }

}

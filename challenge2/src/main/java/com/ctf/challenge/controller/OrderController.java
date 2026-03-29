package com.ctf.challenge.controller;

import com.ctf.challenge.model.Order;
import com.ctf.challenge.service.OrderService;
import jakarta.servlet.http.HttpSession;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;

@Controller
public class OrderController {

    @Autowired
    private OrderService orderService;

    @GetMapping("/orders")
    public String viewOrders(HttpSession session, Model model) {
        Integer userId = (Integer) session.getAttribute("USER_ID");
        if (userId == null) {
            return "redirect:/login";
        }

        List<Order> orders = orderService.getUserOrders(userId);
        model.addAttribute("orders", orders);
        return "user/orders";
    }

    @GetMapping("/orders/{id}")
    public String trackOrder(@PathVariable Integer id, HttpSession session, Model model) {
        Integer userId = (Integer) session.getAttribute("USER_ID");
        if (userId == null) {
            return "redirect:/login";
        }

        Order order = orderService.getOrderDetails(id);
        if (order == null || !order.getUserId().equals(userId)) {
            return "redirect:/"; // Unauthorized or not found
        }

        model.addAttribute("order", order);
        model.addAttribute("items", orderService.getOrderItems(id));
        return "user/track";
    }
}

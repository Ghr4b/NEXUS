package com.ctf.challenge.controller;

import com.ctf.challenge.service.OrderService;
import com.ctf.challenge.service.ProductService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

@Controller
@RequestMapping("/admin")
public class AdminController {

    @Autowired
    private ProductService productService;

    @Autowired
    private OrderService orderService;

    @GetMapping
    public String dashboard(Model model) {
        model.addAttribute("products", productService.getAllProducts());
        model.addAttribute("orders", orderService.getAllOrders());
        return "admin/dashboard";
    }

    @PostMapping("/orders/update")
    public String updateOrderStatus(@RequestParam Integer orderId, @RequestParam String status) {
        orderService.updateStatus(orderId, status);
        return "redirect:/admin";
    }
}

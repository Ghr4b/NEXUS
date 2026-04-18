package com.ctf.challenge.controller;

import com.ctf.challenge.model.Product;
import com.ctf.challenge.model.PromoCode;
import com.ctf.challenge.repository.PromoRepository;
import com.ctf.challenge.service.OrderService;
import com.ctf.challenge.service.ProductService;
import jakarta.servlet.http.HttpSession;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Controller
public class CartController {

    @Autowired
    private ProductService productService;

    @Autowired
    private PromoRepository promoRepository;

    @Autowired
    private OrderService orderService;

    private static final String PROMO_REGEX = "^CYBERSHOP-(?:[A-Z0-9]++-)+[A-Z0-9]+$";
    private static final Pattern PROMO_PATTERN = Pattern.compile(PROMO_REGEX);

    @GetMapping("/cart")
    public String viewCart(HttpSession session, Model model) {
        Map<Integer, Integer> cart = getCart(session);
        model.addAttribute("cartItems", buildCartDetails(cart));
        model.addAttribute("discount", getDiscount(session));
        return "cart/view";
    }

    @PostMapping("/cart/add")
    public String addToCart(@RequestParam Integer productId, @RequestParam Integer quantity, HttpSession session) {
        Map<Integer, Integer> cart = getCart(session);
        cart.put(productId, cart.getOrDefault(productId, 0) + quantity);
        session.setAttribute("cart", cart);
        return "redirect:/cart";
    }

    @PostMapping("/cart/apply-promo")
    public String applyPromoCode(@RequestParam("promoCode") String promoCode, HttpSession session, Model model) {
        Matcher matcher = PROMO_PATTERN.matcher(promoCode);

        if (!matcher.matches()) {
            model.addAttribute("error", "Invalid promo code format. Must follow CYBERSHOP-XXX-XXX format.");
            return viewCart(session, model);
        }

        PromoCode promo = promoRepository.findByCode(promoCode);
        if (promo != null && promo.getIsActive()) {
            Integer userId = (Integer) session.getAttribute("USER_ID");
            if (promo.getUserId() != null && !promo.getUserId().equals(userId)) {
                model.addAttribute("error", "This promo code is not applicable to you.");
            } else {
                session.setAttribute("promoId", promo.getId());
                session.setAttribute("discount", promo.getDiscountPercentage());
                model.addAttribute("message", "Promo code applied successfully!");
            }
        } else {
            model.addAttribute("error", "Promo code is invalid or expired.");
        }
        return viewCart(session, model);
    }

    @PostMapping("/cart/checkout")
    public String checkout(HttpSession session, Model model) {
        Integer userId = (Integer) session.getAttribute("USER_ID");
        if (userId == null) {
            return "redirect:/login";
        }

        Map<Integer, Integer> cart = getCart(session);
        if (cart.isEmpty()) {
            return "redirect:/cart";
        }

        Integer discountPerc = getDiscount(session);
        BigDecimal discountFactor = BigDecimal.valueOf(1)
                .subtract(BigDecimal.valueOf(discountPerc).divide(BigDecimal.valueOf(100)));

        Integer orderId = orderService.checkout(userId, cart, discountFactor);

        session.removeAttribute("cart");
        session.removeAttribute("promoId");
        session.removeAttribute("discount");

        return "redirect:/orders/" + orderId;
    }

    @SuppressWarnings("unchecked")
    private Map<Integer, Integer> getCart(HttpSession session) {
        Map<Integer, Integer> cart = (Map<Integer, Integer>) session.getAttribute("cart");
        return cart == null ? new HashMap<>() : cart;
    }

    private Integer getDiscount(HttpSession session) {
        Integer discount = (Integer) session.getAttribute("discount");
        return discount == null ? 0 : discount;
    }

    private List<Map<String, Object>> buildCartDetails(Map<Integer, Integer> cart) {
        List<Map<String, Object>> details = new ArrayList<>();
        if (cart != null) {
            cart.forEach((id, qty) -> {
                Product p = productService.getProduct(id);
                if (p != null) {
                    Map<String, Object> item = new HashMap<>();
                    item.put("product", p);
                    item.put("quantity", qty);
                    item.put("subtotal", p.getPrice().multiply(BigDecimal.valueOf(qty)));
                    details.add(item);
                }
            });
        }
        return details;
    }
}

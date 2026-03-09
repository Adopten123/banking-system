use chrono::{Datelike, Utc};
use rand::Rng;

pub struct CardGenerator;

impl CardGenerator {
    pub fn generate_pan(prefix: &str) -> String {
        let mut rng = rand::thread_rng();
        let mut pan = String::from(prefix);

        while pan.len() < 15 {
            let digit: u8 = rng.gen_range(0..=9);
            pan.push_str(&digit.to_string());
        }

        let check_digit = Self::calculate_luhn_check_digit(&pan);
        pan.push_str(&check_digit.to_string());

        pan
    }

    fn calculate_luhn_check_digit(number_without_check: &str) -> u32 {
        let mut sum = 0;
        let mut double = true;

        for ch in number_without_check.chars().rev() {
            if let Some(digit) = ch.to_digit(10) {
                if double {
                    let mut doubled_digit = digit * 2;
                    if doubled_digit > 9 {
                        doubled_digit -= 9;
                    }
                    sum += doubled_digit;
                } else {
                    sum += digit;
                }
                double = !double;
            }
        }

        (10 - (sum % 10)) % 10
    }

    pub fn generate_cvv() -> String {
        let mut rng = rand::thread_rng();
        let cvv: u16 = rng.gen_range(0..=999);
        format!("{:03}", cvv)
    }

    pub fn generate_expiry() -> (i32, i32) {
        let now = Utc::now();
        let month = now.month() as i32;
        let year = now.year() + 3;
        (month, year)
    }

    pub fn mask_pan(pan: &str) -> String {
        if pan.len() < 16 {
            return pan.to_string();
        }
        let first_four = &pan[0..4];
        let last_four = &pan[12..16];
        format!("{} **** **** {}", first_four, last_four)
    }
}
import 'package:flutter/material.dart';
import 'login_page.dart';

class Globaldata {
  static String JWToken = "";
  static String ErrorServerConnection = "Could not find the server\nPlease try again later";
  static BuildContext myContext = myContext;
  static Widget lastPage = LoginPage();
  static int firstService = -1;
  static int secondService = -1;

  static List<String> nameList = [
    "Github",
    "Youtube",
    "Spotify",
  ];
}

import 'package:flutter/material.dart';

import "../globalData.dart";

class AllServicesList {
  static List<String> allServiceNameList = Globaldata.actionServicesName + Globaldata.reactionServicesName;
  static List<String> allServiceColorList = Globaldata.actionServicesColor + Globaldata.reactionServicesColor;
  static List<String> allServiceIconsList = Globaldata.actionServicesIcons + Globaldata.reactionServicesIcons;
}

class SelectedCardList {
  static int index = 0;
  static List<Widget> cardList = [];
}
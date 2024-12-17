import 'package:flutter/material.dart';

import "../globalData.dart";

enum SelectedCardListEnum {
    All,
    Workflows,
    Services,
}

class AllServicesList {
  static List<String> allServiceNameList = Globaldata.actionServicesName + Globaldata.reactionServicesName;
  static List<String> allServiceColorList = Globaldata.actionServicesColor + Globaldata.reactionServicesColor;
  static List<String> allServiceIconsList = Globaldata.actionServicesIcons + Globaldata.reactionServicesIcons;
}

class SelectedCardList {
  static SelectedCardListEnum cardType = SelectedCardListEnum.All;
  static List<Widget> cardList = [];
}